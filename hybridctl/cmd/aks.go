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

	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "Can't open the file")
	}

	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)

	if err != nil {
		err = errors.Wrap(err, "Can't unmarshal file")
	}

	return &contents, err
}

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
					SubscriptionId:    "",
					ResourceGroupName: resourceGroupName,
					ResourceName:      clusterName,
					ApiVersion:        "",
				}
				aksStart(EksAPIParameter)

			case "stop":
				resourceGroupName, _ := cmd.Flags().GetString("resource-group")
				clusterName, _ := cmd.Flags().GetString("name")
				EksAPIParameter := util.EksAPIParameter{
					SubscriptionId:    "",
					ResourceGroupName: resourceGroupName,
					ResourceName:      clusterName,
					ApiVersion:        "",
				}
				aksStop(EksAPIParameter)
			}
			// default:
			// 	fmt.Println("Run 'hybridctl aks --help' to view all commands")
			// }
		}
	},
}

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		EksAPIParameter := util.EksAPIParameter{
			SubscriptionId:    "",
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
			ApiVersion:        "",
		}
		aksStart(EksAPIParameter)

	},
}

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: ` 

	`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		EksAPIParameter := util.EksAPIParameter{
			SubscriptionId:    "",
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
			ApiVersion:        "",
		}
		aksStop(EksAPIParameter)

	},
}

func aksStart(p util.EksAPIParameter) {
	// AzureAuth(p.SubscriptionId)
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

func aksStop(p util.EksAPIParameter) {
	// AzureAuth(p.SubscriptionId)
	httpPostUrl := "http://localhost:8080/aksStop"
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
	aksCmd.AddCommand(StartCmd)
	aksCmd.AddCommand(StopCmd)
	aksCmd.PersistentFlags().StringP("resource-group", "g", "", "resourceGroup name")
	aksCmd.PersistentFlags().StringP("name", "n", "", "clustername")
	aksCmd.MarkPersistentFlagRequired("resource-group")
	aksCmd.MarkPersistentFlagRequired("name")
}
