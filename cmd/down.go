/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/Hawkbawk/falcon/lib/dnsmasq"
	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/Hawkbawk/falcon/lib/networking"
	"github.com/Hawkbawk/falcon/lib/proxy"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Restores networking to normal and stops the proxy and dnsmasq containers",
	Long: `Running the down command will restore all machine networking to normal. It will also
	stop and remove the proxy and dnsmasq container.`,
	Run: func(cmd *cobra.Command, args []string) {
		networking.Clean()
		client, err := docker.NewDockerClient()
		if err != nil {
			logger.LogError("Unable to connect to Docker server due to the following error:\n%v", err)
			os.Exit(1)
		}

		logger.LogInfo("Stopping the dnsmasq container...")
		if err := dnsmasq.Stop(client); err != nil {
			logger.LogError("Unable to stop the dnsmasq container:\n%v", err)
			os.Exit(1)
		}
		logger.LogInfo("Stopping the falcon proxy container...")
		if err := proxy.Stop(client); err != nil {
			logger.LogError("Unable to stop the proxy container:\n%v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
