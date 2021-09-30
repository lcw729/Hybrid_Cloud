package cmd

import (
	"Hybrid_Cluster/hybridctl/util"

	"github.com/spf13/cobra"
)

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
