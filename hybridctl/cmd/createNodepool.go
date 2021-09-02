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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// createNodepoolCmd represents the createNodepool command
var createNodepoolCmd = &cobra.Command{
	Use:   "createNodepool",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("createNodepool called")
		if len(args) == 0 {
			fmt.Println("Run 'hybridctl createNodepool --help' to view all commands")
		} else if args[0] == "gke" {
			if args[1] == "" {
				fmt.Println("Run 'hybridctl createNodepool --help' to view all commands")
			} else {
				fmt.Println("kubernetes engine Name: ", args[0])
				fmt.Printf("Cluster Name: %s\n", args[1])
				createNodepool_gke(args[1])
			}
		}
	},
}

type TF_NodePool struct {
	NodePool_Resource *NodePool_Resource `json:"resource"`
}

type NodePool_Resource struct {
	Google_container_node_pool *map[string]Node_pool_type `json:"google_container_node_pool"`
}

type Node_pool_type struct {
	Name        string       `json:"name"`
	Location    string       `json:"location"`
	Cluster     string       `json:"cluster"`
	Node_count  int          `json:"node_count"`
	Node_config *Node_config `json:"node_config"`
}

type Labels struct {
	Env string `json:"env"`
}
type Node_config struct {
	Oauth_scopes []string  `json:"oauth_scopes"`
	Labels       *Labels   `json:"labels"`
	Machine_type string    `json:"machine_type"`
	Tags         []string  `json:"tags"`
	Metadata     *Metadata `json:"metadata"`
}

type Metadata struct {
	Disable_legacy_endpoints string `json:"disable-legacy-endpoints"`
}

func createNodepool_gke(clusterName string) {
	cluster := "cluster"
	num := 1

	send_js_nodePool := TF_NodePool{
		NodePool_Resource: &NodePool_Resource{
			Google_container_node_pool: &map[string]Node_pool_type{
				clusterName + "_nodes": {
					Name:       "${google_container_cluster." + clusterName + ".name}-node-pool",
					Location:   "us-central1-a",
					Cluster:    "${google_container_cluster." + clusterName + ".name}",
					Node_count: num,
					Node_config: &Node_config{
						Labels: &Labels{
							Env: "keti-container",
						},
						Metadata: &Metadata{
							Disable_legacy_endpoints: "true",
						},
						Tags:         []string{"gke-node", "keti-container-gke"},
						Machine_type: "n1-standard-1",
						Oauth_scopes: []string{"https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/monitoring"},
					},
				},
			},
		},
	}

	send, err := json.MarshalIndent(send_js_nodePool, "", " ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("/root/Hybrid_Cluster/terraform/gke/"+cluster+"/"+clusterName+"nodePool"+".tf.json", []byte(string(send)), os.FileMode(0644))
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("terraform", "apply", "-auto-approve")
	// cmd := exec.Command("terraform", "plan")
	cmd.Dir = "../terraform/gke/" + cluster

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}
}

func init() {
	RootCmd.AddCommand(createNodepoolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createNodepoolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createNodepoolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
