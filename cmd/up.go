/*
Copyright © 2021 Ryan Hawkins ryanlarryhawkins@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/Hawkbawk/falcon/lib/dnsmasq"
	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/Hawkbawk/falcon/lib/networking"
	"github.com/Hawkbawk/falcon/lib/proxy"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Sets up networking and starts the dnsmasq and proxy container",
	Long: `falcon up sets up your local networking to point all requests to *.docker to resolve
to localhost:80, and then starts the dnsmasq and proxy container.
The proxy container (running Traefik) then takes these requests and acts as
a reverse-proxy, determining to which container the request should go to.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := networking.Configure(); err != nil {
			logger.LogError("Couldn't configure networking:\n%v", err)
		}
		client, err := docker.NewDockerClient()
		if err != nil {
			logger.LogError("Unable to connect to the Docker server:\n%v", err)
		}

		logger.LogInfo("Starting the dnsmasq container...")
		if err := dnsmasq.Start(client); err != nil {
			logger.LogError("Unable to start the dnsmasq container: \n%v", err)
		}

		logger.LogInfo("Starting the proxy container...")
		if err := proxy.Start(client); err != nil {
			logger.LogError("Unable to start the dnsmasq container:\n%v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
