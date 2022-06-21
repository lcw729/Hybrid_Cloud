package cmd

import (
	"Hybrid_Cloud/hybridctl/util"
	"fmt"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var file_name string

// DeleteCmd represents the Delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `hybridctl delete deployment <name> -n <namespace> 
	hybridctl delete -f <filename> `,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		file_name = util.Option_file
		if file_name == "" {
			if len(args) < 2 {
				fmt.Println(cmd.Help())
			} else {
				LINK := "/resources"
				namespace, _ := cmd.Flags().GetString("namespace")
				if namespace == "" {
					namespace = "default"
				}
				LINK += "/namespaces/" + namespace

				util.Option_Resource = args[0]
				util.Option_Name = args[1]
				LINK += "/deployments/" + util.Option_Name

				_, err := util.GetResponseBody("DELETE", LINK, nil)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			DeleteResource()
		}
	},
}

func DeleteResource() {

	yaml, err := ReadFile()
	if err != nil {
		println(err)
		return
	}

	obj, gvk, err := GetObject(yaml)
	if err != nil {
		println(err)
		return
	}

	RequestDeleteResource(obj, gvk)
}

func RequestDeleteResource(obj runtime.Object, gvk *schema.GroupVersionKind) ([]byte, error) {

	LINK := "/resources"
	// check context flag
	//	flag_context := util.Option_context
	// var target_cluster string
	// var resource Resource

	// config, _ := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
	// 	&clientcmd.ClientConfigLoadingRules{ExplicitPath: "/root/.kube/config"},
	// 	&clientcmd.ConfigOverrides{
	// 		CurrentContext: "",
	// 	}).RawConfig()

	// if flag_context == "" {
	// 	target_cluster = ""
	// } else {
	// 	target_cluster = flag_context
	// }

	// match obj kind
	switch gvk.Kind {
	case "Deployment":
		real_resource := obj.(*appsv1.Deployment)
		namespace := real_resource.Namespace
		if namespace == "" {
			namespace = "default"
		}
		LINK += "/namespaces/" + namespace + "/deployments/" + real_resource.Name
	}

	fmt.Println(LINK)
	bytes, err := util.GetResponseBody("DELETE", LINK, nil)
	if err != nil {
		fmt.Println(err)
	}

	return bytes, err
}

func init() {
	RootCmd.AddCommand(DeleteCmd)
	DeleteCmd.Flags().StringVarP(&util.Option_file, "file", "f", "", "FILENAME")
	DeleteCmd.MarkFlagRequired("file")
	DeleteCmd.Flags().StringVarP(&util.Option_context, "context", "", "", "CLUSTERNAME")
	DeleteCmd.Flags().StringP("namespace", "n", "default", "enter the namespace")
}
