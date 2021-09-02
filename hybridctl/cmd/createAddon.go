// Copyright © 2021 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"Hybrid_Cluster/hybridctl/util"
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/cobra"
)

var createAddonInput eks.CreateAddonInput

// createAddonCmd represents the createAddon command
var createAddonCmd = &cobra.Command{
	Use:   "create-addon",
	Short: "A brief description of your command",
	Long: `	
	- create-addon
		hybridctl create-addon <clusterName> <addonName> 

	- flags
		[--addon-version <value>]
		[--service-account-role-arn <value>]
		[--resolve-conflicts <value>]
		[--client-request-token <value>]
		[--tags <value>]

	- platform
		- eks (elastic kubernetes service)`,

	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		if len(args) == 0 || len(args) == 1 {
			fmt.Println("Run 'hybridctl join --help' to view all commands")
		} else if args[1] == "" {
			fmt.Println("Run 'hybridctl join --help' to view all commands")
		} else {
			addonVersion, _ := cmd.Flags().GetString("addon-version")
			serviceAccountRoleArn, _ := cmd.Flags().GetString("service-account-role-arn")
			resolveConflicts, _ := cmd.Flags().GetString("resolve-conflicts")
			clientRequestToken, _ := cmd.Flags().GetString("client-request-token")
			createAddonInput.ClusterName = &args[0]
			createAddonInput.AddonName = &args[1]
			if addonVersion != "" {
				createAddonInput.AddonVersion = &addonVersion
			}
			if serviceAccountRoleArn != "" {
				createAddonInput.ServiceAccountRoleArn = &serviceAccountRoleArn
			}
			if resolveConflicts != "" {
				createAddonInput.ResolveConflicts = &resolveConflicts
			}
			if clientRequestToken != "" {
				createAddonInput.ClientRequestToken = &clientRequestToken
			}
			tags, _ := cmd.Flags().GetString("tags")
			var tagsMap map[string]*string
			if tags != "" {
				util.UnmarshalJsonFile(tags, &tagsMap)
				createAddonInput.Tags = tagsMap
			}
			// createAddonInput.Tags = tags
			createAddon(createAddonInput)
		}
	},
}

func createAddon(createAddonInput eks.CreateAddonInput) {
	httpPostUrl := "http://localhost:8080/createAddon"
	var output eks.CreateAddonOutput
	util.GetJson(httpPostUrl, createAddonInput, &output)
	if output == (eks.CreateAddonOutput{}) {
		fmt.Println("Addon already exists.")
	} else {
		fmt.Printf("%+v\n", output)
	}
}

func init() {
	RootCmd.AddCommand(createAddonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createAddonCmd.Flags().StringP("addon-version", "", "", "enter addon version")
	createAddonCmd.Flags().StringP("service-account-role-arn", "", "", "enter service account rolearn")
	createAddonCmd.Flags().StringP("resolve-conflicts", "", "", "enter addon version")
	createAddonCmd.Flags().StringP("client-request-token", "", "", "enter client request token")
	createAddonCmd.Flags().StringP("tags", "", "", "enter your tags Jsonfile name")
}
