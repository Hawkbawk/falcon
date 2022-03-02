/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/Hawkbawk/falcon/lib/proxy"
	"github.com/spf13/cobra"
)

// tlsCmd represents the tls command
var tlsCmd = &cobra.Command{
	Use:   "tls <domain_name>",
	Short: "The tls command lets you enable TLS for a domain of your choosing",
	Long: `You can use the tls command to enable TLS/HTTPS for
the domain of your choosing. If your application needs HTTPS
to work properly, this is the way to do it. Just run the tls command
on the domain that matches your application and then enable TLS for your
container by putting a label that matches the format:
"traefik.http.routers.<your_router_name>.tls=true"
`,
	ValidArgs: []string{"domain_name"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logger.LogError("You must specify a hostname and only a hostname!")
		}

		if err := proxy.EnableTlsForHost(args[0]); err != nil {
			logger.LogError("Unable to enable TLS for the specified host:\n%v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tlsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tlsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tlsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
