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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	mappingTable "Hybrid_Cluster/hcp-apiserver/pkg/converter"

	cobrautil "Hybrid_Cluster/hybridctl/util"

	hcpclusterapis "Hybrid_Cluster/pkg/apis/hcpcluster/v1alpha1"
	hcpclusterv1alpha1 "Hybrid_Cluster/pkg/client/hcpcluster/v1alpha1/clientset/versioned"

	util "Hybrid_Cluster/util"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var checkAKS, checkEKS, checkGKE = false, false, false
var master_config, _ = cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
var master_client = kubernetes.NewForConfigOrDie(master_config)

type Cli struct {
	PlatformName string
	ClusterName  string
}

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "A brief description of your command",
	Long: ` 
NAME 
	hybridctl join PLATFORM CLUSTER
	hybridctl join register PLATFORM

DESCRIPTION
	
	>> cluster join PLATFORM CLUSTER <<


	PLATFORM means the Kubernetes platform of the cluster to join.
	The types of platforms offered are as follows.

	- aks   azure kubernetes service
	- eks   elastic kubernetes service
	- gke   google kuberntes engine

	* PLATFORM mut be written in LOWERCASE letters

	CLUSTER means the name of the cluster on the specified platform.

	>> hybridctl join register PLATFORM <<

	* This command registers the cluster you want to manage, 
	For each platform, you must fill in the information below.
	Please refer to the INFO section

	PLATFORM means the Kubernetes platform of the cluster to join.
	The types of platforms offered are as follows.

	- aks   azure kubernetes service
	- eks   elastic kubernetes service
	- gke   google kuberntes engine

	[INFO]

		GKE 
		- projectid    the ID of GKE cloud project to use. 
		- clustername  the name of the cluster on the specified platform.
		- region       choose Google Compute Zone from 1 to 85.

	`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		if len(args) == 0 {
		} else {
			switch args[0] {
			case "aks":
				fallthrough
			case "eks":
				fallthrough
			case "gke":
				fmt.Println("kubernetes engine Name : ", args[0])
				fmt.Printf("Cluster Name : %s\n", args[1])
				cli := mappingTable.ClusterInfo{
					PlatformName: args[0],
					ClusterName:  args[1]}
				join(cli)
			case "register":
				platform := args[1]
				var clustername string
				fmt.Printf("please enter your cluster name: ")
				fmt.Scanln(&clustername)
				createPlatformNamespace()
				// path := "/root/Go/src/Hybrid_Cluster/hybridctl/cmd/kubeconfig.sh"
				// CmdExecsh(path)
				// config := util.KubeConfig{}
				switch platform {
				case "aks":
					var resourcegroup string
					fmt.Printf("please enter your resourceGroup: ")
					fmt.Scanln(&resourcegroup)
					// CmdExec("az aks get-credentials --resource-group " + resourcegroup + " --name " + clustername + " --file ~/.kube/kubeconfig")
					args := []string{"/root/Go/src/Hybrid_Cluster/hybridctl/cmd/join-register.sh", platform, clustername, resourcegroup}
					CmdExecsh("/root/Go/src/Hybrid_Cluster/hybridctl/cmd/join-register.sh", args)
					CreateHCPCluster(clustername, platform)
				case "eks":
					// CmdExec("aws eks update-kubeconfig --name " + clustername)
					args = []string{platform, clustername}
					CmdExecsh("/root/Go/src/Hybrid_Cluster/hybridctl/cmd/join-register.sh", args)
					CreateHCPCluster(clustername, platform)
				case "gke":
					// CmdExec("gcloud container clusters get-credentials " + clustername)
					args = []string{platform, clustername}
					CmdExecsh("/root/Go/src/Hybrid_Cluster/hybridctl/cmd/join-register.sh", args)
					CreateHCPCluster(clustername, platform)
				default:
					return
				}
				// path = "/root/Go/src/Hybrid_Cluster/hybridctl/cmd/config.sh"
				// CmdExecsh(path)

				// clusterRegisterClientSet, err := clusterRegisterv1alpha1.NewForConfig(master_config)
				// var region string
				// var clustername string
				// if mcerr != nil {
				// 	log.Println(err)
				// }
				// createPlatformNamespace()
				/*
					switch args[1] {
					case "aks":
						var resourcegroup string
						fmt.Printf("please enter your cluster region: ")
						fmt.Scanln(&region)
						fmt.Printf("Enter resourcegroup : ")
						fmt.Scanln(&resourcegroup)
						fmt.Printf("clustername : ")
						fmt.Scanln(&clustername)
						newclusterRegister := &resourcev1alpha1.ClusterRegister{
							TypeMeta: metav1.TypeMeta{
								Kind:       "ClusterRegister",
								APIVersion: "hcp.k8s.io/v1alpha1",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      clustername,
								Namespace: "aks",
							},
							Spec: resourcev1alpha1.ClusterRegisterSpec{
								Name:          clustername,
								Region:        region,
								Platform:      "aks",
								Resourcegroup: resourcegroup,
							},
						}

						_, err = clusterRegisterClientSet.ClusterRegisters("aks").Create(context.TODO(), newclusterRegister, metav1.CreateOptions{})

						if err != nil {
							log.Println(err)
						}
					case "eks":
						fmt.Printf("Enter cluster region : ")
						fmt.Scanln(&region)
						fmt.Printf("Enter clustername : ")
						fmt.Scanln(&clustername)
						newclusterRegister := &resourcev1alpha1.ClusterRegister{
							TypeMeta: metav1.TypeMeta{
								Kind:       "ClusterRegister",
								APIVersion: "hcp.k8s.io/v1alpha1",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      clustername,
								Namespace: "eks",
							},
							Spec: resourcev1alpha1.ClusterRegisterSpec{
								Name:     clustername,
								Region:   region,
								Platform: "eks",
							},
						}
						_, err = clusterRegisterClientSet.ClusterRegisters("eks").Create(context.TODO(), newclusterRegister, metav1.CreateOptions{})

						if err != nil {
							log.Println(err)
						}
					case "gke":
						var projectid string
						fmt.Printf("Please enter cloud projectID: ")
						fmt.Scanln(&projectid)
					LABEL:
						fmt.Printf("Please choice Google Compute Engine zone where your cluster exists?\n")
						for i := 1; i <= 50; i++ {
							fmt.Printf("[%d] %s\n", i, cobrautil.GKEregion[i])
						}
						fmt.Printf("Too many options [85]. Enter \"list\" at prompt to print choices fully.\nPlease enter numeric choice or text value (must exactly match list item):")
						fmt.Scanln(&region)
						if region == "list" {
							for i := 0; i < 85; i++ {
								fmt.Printf("[%d] %s\n", i, cobrautil.GKEregion[i])
							}
						}
						num, err := strconv.Atoi(region)
						if err != nil {
							log.Println(err)
						}
						if region != "list" && reflect.TypeOf(region).Kind() == reflect.Int {
							if 1 > num || num > 85 {
								fmt.Printf("Please enter numeric choice \n")
								goto LABEL
							}
						}
						fmt.Printf("Please enter your the name of the cluster to register: ")
						fmt.Scanln(&clustername)
						// gke_cr := strings.Replace(projectid+"-"+cobrautil.GKEregion[num]+"-"+clustername+"-hcp", " ", "", -1)
						newclusterRegister := &resourcev1alpha1.ClusterRegister{
							TypeMeta: metav1.TypeMeta{
								Kind:       "ClusterRegister",
								APIVersion: "hcp.k8s.io/v1alpha1",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      clustername,
								Namespace: "gke",
							},
							Spec: resourcev1alpha1.ClusterRegisterSpec{
								Name:      clustername,
								Region:    cobrautil.GKEregion[num],
								Platform:  "gke",
								ProjectId: projectid,
							},
						}

						_, err = clusterRegisterClientSet.ClusterRegisters("gke").Create(context.TODO(), newclusterRegister, metav1.CreateOptions{})

						if err != nil {
							log.Println(err)
						}
					default:
						fmt.Println("Run 'hybridctl join --help' to view all commands")
					}
				*/
			default:
				fmt.Println("Run 'hybridctl join --help' to view all commands")
			}
		}
	},
}

func CmdExec(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Env = append(cmd.Env, "KUBECONFIG=~/.kube/kubeconfig")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(output), err
}

func CmdExecsh(path string, args []string) (string, error) {
	cmd := exec.Command("/bin/sh", args...)
	cmd.Args = args
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(output), err
}

func CreateHCPCluster(clustername string, platform string) {
	hcp_cluster, err := hcpclusterv1alpha1.NewForConfig(master_config)
	if err != nil {
		log.Println(err)
	}
	var config util.KubeConfig

	kubeconfig, err := json.Marshal(config)
	if err != nil {
		log.Println(err)
	}
	cluster := hcpclusterapis.HCPCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HCPCluster",
			APIVersion: "hcp.k8s.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clustername,
			Namespace: platform,
		},
		Spec: hcpclusterapis.HCPClusterSpec{
			ClusterPlatform: platform,
			KubeconfigInfo:  kubeconfig,
			JoinStatus:      "UNJOIN",
		},
	}
	newhcpcluster, err := hcp_cluster.HcpV1alpha1().HCPClusters(platform).Create(context.TODO(), &cluster, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("success to register %s in %s", newhcpcluster.Name, newhcpcluster.Namespace)
	}
}
func join(info mappingTable.ClusterInfo) {
	httpPostUrl := "http://localhost:8080/join"
	jsonData, _ := json.Marshal(&info)

	buff := bytes.NewBuffer(jsonData)
	request, _ := http.NewRequest("POST", httpPostUrl, buff)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	// handler.Join(info)
}

func createPlatformNamespace() {

	namespaceList, _ := master_client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	for i := range namespaceList.Items {
		if checkAKS && checkEKS && checkGKE {
			break
		}
		switch namespaceList.Items[i].Name {
		case "aks":
			checkAKS = true
			continue
		case "eks":
			checkEKS = true
			continue
		case "gke":
			checkGKE = true
			continue
		default:
			continue
		}
	}
	checkAndCreateNamespace(checkAKS, "aks")
	checkAndCreateNamespace(checkEKS, "eks")
	checkAndCreateNamespace(checkGKE, "gke")
}

func checkAndCreateNamespace(PlatformCheck bool, platformName string) {
	if !PlatformCheck {
		Namespace := corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: platformName,
			},
		}
		master_client.CoreV1().Namespaces().Create(context.TODO(), &Namespace, metav1.CreateOptions{})
	}
}

func init() {
	RootCmd.AddCommand(joinCmd)
}
