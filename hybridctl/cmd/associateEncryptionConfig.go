package cmd

import (
	"Hybrid_Cluster/hybridctl/util"
	cobrautil "Hybrid_Cluster/hybridctl/util"
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/spf13/cobra"
)

var associateEncryptionConfigInput eks.AssociateEncryptionConfigInput

// AssociateIdentityProvicerConfigCmd represents the AssociateIdentityProvicerConfig command
var associateEncryptionConfigCmd = &cobra.Command{
	Use:   "associate-encryption-config",
	Short: "A brief description of your command",
	Long: `	
	- associate-encryption-config
		hybridctl associate-encryption-config <clusterName> --encryption-config <jsonfile>

	- platform
		- eks (elastic kubernetes service)`,

	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		if len(args) == 0 {
			fmt.Println("Run 'hybridctl associate-encryption-config --help' to view all commands")
		} else if args[0] == "" {
			fmt.Println("Run 'hybridctl associate-encryption-config --help' to view all commands")
		} else {
			associateEncryptionConfigInput.ClusterName = &args[0]

			// json parsing
			jsonFileName, _ := cmd.Flags().GetString("encryption-config")
			var encryptionConfig []*eks.EncryptionConfig
			util.UnmarshalJsonFile(jsonFileName, encryptionConfig)
			associateEncryptionConfigInput.EncryptionConfig = encryptionConfig

			clientRequestToken, _ := cmd.Flags().GetString("client-request-token")
			if clientRequestToken != "" {
				associateEncryptionConfigInput.ClientRequestToken = &clientRequestToken
			}
			AssociateEncryptionConfig(associateEncryptionConfigInput)
		}
	},
}

func AssociateEncryptionConfig(AssociateEncryptionConfigInput eks.AssociateEncryptionConfigInput) {
	httpPostUrl := "http://localhost:8080/associateEncryptionConfig"
	var output eks.AssociateEncryptionConfigOutput
	cobrautil.GetJson(httpPostUrl, AssociateEncryptionConfigInput, &output)
	fmt.Printf("%+v\n", output)
}

func init() {
	EksCmd.AddCommand(associateEncryptionConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// AssociateIdentityProvicerConfigCmd.Flags().StringVarP(&cobrautil.Option_context, "context", "c", "", "input a option")
	// associateEncryptionConfigCmd.Flags().StringP("tags", "", "", "enter tags")
	associateEncryptionConfigCmd.Flags().StringP("encryption-config", "", "", "enter your encryption-config Jsonfile name")
	associateEncryptionConfigCmd.Flags().StringP("client-request-token", "", "", "enter client request token")
}
