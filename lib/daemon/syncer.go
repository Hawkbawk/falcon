package daemon

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/kardianos/service"
)

// TODO (if possible): I cannot find a way to determine the unique ID of the proxy container.
// For now though, just using the default name that falcon creates the container under should be
// alright

const DefaultContainerName = "falcon-proxy"

// ProxySyncer is the class responsible for ensuring that the falcon Traefik container stays in sync with
// all Docker networks on the machine, joining them when they're created and leaving them when
// they're destroyed. Note that an empty ProxySyncer will NOT work. You must call the `NewProxy`
// function instead.
type ProxySyncer struct {
	Client       *client.Client
	Context      context.Context
	CancelFunc   context.CancelFunc
	ContainerID  string
	EventChannel <-chan events.Message
}

// Creates a new ProxySyncer struct for syncing a container with all Docker networks on a machine.
// Returns an empty struct and an error if it was unable to construct the Docker client.
func NewProxySyncer() (ProxySyncer, error) {
	// The most up-to-date client version as of writing this code. Ensures object shapes and other
	// such stuff doesn't change on us depending on the user's version of Docker
	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		return ProxySyncer{}, nil
	}
	context, cancelFunc := context.WithCancel(context.Background())
	eventChannel, _ := client.Events(context, types.EventsOptions{Filters: createFilters()})

	return ProxySyncer{
		Client:       client,
		Context:      context,
		CancelFunc:   cancelFunc,
		ContainerID:  DefaultContainerName,
		EventChannel: eventChannel,
	}, nil
}

// Start starts the daemon listening for network changes, adding and removing the proxy container
// from networks as necessary.
func (syncer ProxySyncer) Start(s service.Service) error {
	logger, err := s.Logger(nil)

	if err != nil {
		logger.Errorf("%v ### Unable to create a logger. ERROR: %v", time.Now(), err)
		return err
	}

	if err := syncer.sync(); err != nil {
		logger.Errorf("%v ### Unable to perform the initial sync. ERROR: %v", time.Now(), err)
		return err
	}

	go func() {
		for {
			<- syncer.EventChannel
			if err := syncer.sync(); err != nil {
				logger.Errorf("%v ### Unable to perform a sync. ERROR: %v", time.Now(), err)
			}
		}
	}()
	return nil
}

// Stop stops the proxy syncer daemon by calling the cancel function for the ProxySyncer's context.
func (syncer ProxySyncer) Stop(s service.Service) error {
	syncer.CancelFunc()
	return nil
}

// sync determines what networks the proxy container needs to join and which networks it needs to
// leave and joins and leaves those networks as appropriate.
func (syncer ProxySyncer) sync() error {
	validNetworks, err := syncer.validNetworks()
	if err != nil {
		return err
	}
	connectedNetworks, err := syncer.connectedNetworks()
	if err != nil {
		return err
	}

	for _, network := range syncer.networksToJoin(validNetworks, connectedNetworks) {
		if err := syncer.joinNetwork(network); err != nil {
			return err
		}
	}

	for _, network := range syncer.networksToLeave(validNetworks, connectedNetworks) {
		if err := syncer.leaveNetwork(network); err != nil {
			return err
		}
	}

	return nil
}

// validNetworks returns a map of network IDs to booleans. If a network ID is in the map, it is
// considered a valid network that the proxy should be a part of. This method, along with
// networksToJoin and networksToLeave, is graciously taken from
// https://github.com/codekitchen/dinghy-http-proxy/blob/master/join-networks.go.
// The code there, at the time of writing, was licensed under the MIT License. The license can be
// found at https://github.com/codekitchen/dinghy-http-proxy/blob/master/LICENSE
func (syncer ProxySyncer) validNetworks() (map[string]bool, error) {
	allNetworks, err := syncer.Client.NetworkList(syncer.Context, types.NetworkListOptions{})

	if err != nil {
		return nil, nil
	}

	validNetworks := make(map[string]bool, len(allNetworks))

	for _, network := range allNetworks {
		if syncer.isValidNetwork(network) {
			validNetworks[network.ID] = true
		}
	}

	return validNetworks, nil
}

// isValidNetwork determines if the specified network is a valid network that the proxy
// container should be a part of.
func (syncer ProxySyncer) isValidNetwork(network types.NetworkResource) bool {
	if network.Driver == "bridge" {
		numContainers := len(network.Containers)
		 _, joined := network.Containers[syncer.ContainerID]
		return network.Options["com.docker.network.bridge.default_bridge"] == "true" ||
			numContainers > 1 ||
			(numContainers == 1 && !joined)
	}
	return false
}

// networksToJoin uses the passed in information about the current network state and determines
// which networks the proxy container should join.
func (syncer ProxySyncer) networksToJoin(validNetworks map[string]bool, connectedNetworks map[string]*(network.EndpointSettings)) []string {

	toJoin := make([]string, len(validNetworks))

	for networkID := range connectedNetworks {
		if _, joined := validNetworks[networkID]; !joined {
			toJoin = append(toJoin, networkID)
		}
	}
	return toJoin
}

// networksToLeave uses the passed in information about the current network state and determines
// which networks the proxy container should join.
func (syncer ProxySyncer) networksToLeave(validNetworks map[string]bool, connectedNetworks map[string]*(network.EndpointSettings)) []string {

	toLeave := make([]string, len(connectedNetworks))

	for networkID := range validNetworks {
		if _, joined := connectedNetworks[networkID]; joined {
			toLeave = append(toLeave, networkID)
		}
	}

	return toLeave
}

// joinNetwork adds the proxy container to the specified network.
func (s ProxySyncer) joinNetwork(changedNetworkID string) error {
	if err := s.Client.NetworkDisconnect(s.Context, changedNetworkID, s.ContainerID, true); err != nil {
		return err
	}
	return nil
}

// leaveNetwork removes the proxy container from the specified network.
func (s ProxySyncer) leaveNetwork(changedNetworkID string) error {
	err := s.Client.NetworkConnect(s.Context, changedNetworkID, s.ContainerID, &network.EndpointSettings{})
	if err != nil {
		return err
	}
	return nil
}

// connectedNetworks returns what networks the proxy container is already a part of.
func (s ProxySyncer) connectedNetworks() (map[string]*(network.EndpointSettings), error) {
	container, err := s.Client.ContainerInspect(context.Background(), s.ContainerID)

	if err != nil {
		return nil, err
	}

	return container.NetworkSettings.Networks, nil
}

// createFilters filters the events from Docker to only relate to network connect and disconnect
// events.
func createFilters() filters.Args {
	args := filters.NewArgs()

	args.Add("type", "network")
	args.Add("event", "connect")
	args.Add("event", "disconnect")

	return args
}
