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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
)

// aksCmd represents the aks command
var aksCmd = &cobra.Command{
	Use:   "aks",
	Short: "A brief description of your command",
	Long: ` 

	`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		if len(args) < 1 {
			fmt.Println("Run 'hybridctl aks --help' to view all commands")
		} else {
			// switch args[0] {
			// case "aks":
			switch args[0] {
			case "start":
				resourceGroupName, _ := cmd.Flags().GetString("resource-group")
				clusterName, _ := cmd.Flags().GetString("name")
				EksAPIParameter := util.EksAPIParameter{
					SubscriptionId:    "ccfc0c6c-d3c6-4de2-9a6c-c09ca498ff73",
					ResourceGroupName: resourceGroupName,
					ResourceName:      clusterName,
					ApiVersion:        "2021-05-01",
				}
				aksStart(EksAPIParameter)
			}
			// default:
			// 	fmt.Println("Run 'hybridctl aks --help' to view all commands")
			// }
		}
	},
}

func aksStart(p util.EksAPIParameter) {
	AzureAuth(p.SubscriptionId)
	httpPostUrl := "http://localhost:8080/aksStart"
	jsonData, _ := json.Marshal(&p)

	buff := bytes.NewBuffer(jsonData)
	request, _ := http.NewRequest("POST", httpPostUrl, buff)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(bytes))
}

func init() {
	RootCmd.AddCommand(aksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	aksCmd.Flags().StringP("resource-group", "g", "", "resourceGroup name")
	aksCmd.Flags().StringP("name", "n", "", "clustername")
	aksCmd.MarkPersistentFlagRequired("resource-group")
	aksCmd.MarkPersistentFlagRequired("name")
}

func AzureAuth(subscriptionID string) compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Auth: Successful")
		vmClient.Authorizer = authorizer
	}

	return vmClient
}
