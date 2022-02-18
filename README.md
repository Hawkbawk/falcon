# falcon

falcon is your friendly neighborhood local development Docker reverse proxy.
That mouthful of words basically means that falcon lets you access your Docker
containers at friendly domain addresses, like `yourwebapp.docker` instead of
`localhost:<whatever_port_was_free>`. Behind the scenes, falcon uses
[Traefik](https://traefik.io/traefik/) to perform all proxying and some DNS
trickery to get everything working properly.

# Installation

Because falcon is written in Go, the easiest way to install it is to use
Go's module system. After [installing Go](https://go.dev/doc/install),
simply run `go install github.com/Hawkbawk/falcon@latest` to install falcon.
Assuming you've added your GOPATH to your regular PATH, you should be able to
run `falcon` and you should see a friendly message explaining what commands
are available to you. For now, there are only two commands, up, that
start falcon and its requisite services, and down, which stops falcon and
restores your networking configuration to its default state.

# Configuration

Because falcon uses Traefik behind the scense for all proxying, you'll be using
labels to tell falcon the domain (and possibly poprt) that you want to be able to access your container
at. The template below shows you the labels that you'd need to add in order
to access your application at the domain of your choice:

```yaml
- traefik.enable=true # required to make proxying work
- traefik.http.routers.<app_name_here>.rule=Host(`<desired_domain>.docker`) # required to specify domain
- traefik.http.routers.<app_name_here>.loadbalancer.port=80 # optional
```

Note that the port option allows you to specify what port your container is
running your application on. This option is only necessary if your container
exposes multiple ports. If your container only exposes and works on a single
port, Traefik will automatically detect what port to use.

For further reading, see [Traefik's documentation](https://doc.traefik.io/traefik/routing/providers/docker/)
related to routing with Docker
