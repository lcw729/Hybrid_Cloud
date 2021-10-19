package cmd

import (
	"Hybrid_Cluster/hybridctl/util"
	cobrautil "Hybrid_Cluster/hybridctl/util"
	"fmt"

	"github.com/spf13/cobra"
)

//addon
var AddonCmd = &cobra.Command{
	Use:   "addon",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
}

var AKSDisableAddonsCmd = &cobra.Command{
	Use:   "disable",
	Short: "A brief description of your command",
	Long:  `hybridctl aks disable-addons`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		addon, _ := cmd.Flags().GetString("addon")

		AKSAddon := util.AKSAddon{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Addon:             addon,
		}
		addonDisable(AKSAddon)
	},
}

var AKSEnableAddonsCmd = &cobra.Command{
	Use:   "enable",
	Short: "A brief description of your command",
	Long:  `hybridctl aks addon enable`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		addon, _ := cmd.Flags().GetString("addon")

		AKSAddon := util.AKSAddon{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Addon:             addon,
		}
		addonEnable(AKSAddon)
	},
}

var AKSListAddonsCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `hybridctl aks addon list`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		AKSAddon := util.AKSAddon{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
		}
		addonList(AKSAddon)
	},
}

var AKSListAddonsAvailableCmd = &cobra.Command{
	Use:   "list-available",
	Short: "A brief description of your command",
	Long:  `hybridctl aks addon list`,
	Run: func(cmd *cobra.Command, args []string) {
		addonListAvailable()
	},
}

var AKSShowAddonsCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long:  `hybridctl aks addon enable`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		addon, _ := cmd.Flags().GetString("addon")

		AKSAddon := util.AKSAddon{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Addon:             addon,
		}
		addonShow(AKSAddon)
	},
}

var AKSUpdateAddonsCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long:  `hybridctl aks addon enable`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		addon, _ := cmd.Flags().GetString("addon")

		AKSAddon := util.AKSAddon{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Addon:             addon,
		}
		addonUpdate(AKSAddon)
	},
}

//pod-identity
var AKSPodIdentityCmd = &cobra.Command{
	Use:   "pod-identity",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity`,
}

var AKSPIAddCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity add`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		identityResourceID, _ := cmd.Flags().GetString("identity-resource-id")
		namespace, _ := cmd.Flags().GetString("namespace")
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		podIdentityName, _ := cmd.Flags().GetString("name")
		bindingSelector, _ := cmd.Flags().GetString("addon")

		AKSPodIdentity := util.AKSPodIdentity{
			ResourceGroupName:  resourceGroupName,
			ClusterName:        clusterName,
			Namespace:          namespace,
			IdentityResourceID: identityResourceID,
		}
		if podIdentityName != "" {
			AKSPodIdentity.Name = podIdentityName
		}
		if bindingSelector != "" {
			AKSPodIdentity.BindingSelector = bindingSelector
		}
		podIdentityAdd(AKSPodIdentity)
	},
}

var AKSPIDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity delete`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		namespace, _ := cmd.Flags().GetString("namespace")
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		podIdentityName, _ := cmd.Flags().GetString("name")

		AKSPodIdentity := util.AKSPodIdentity{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Namespace:         namespace,
			Name:              podIdentityName,
		}
		podIdentityDelete(AKSPodIdentity)
	},
}

var AKSPIListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity list`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")

		AKSPodIdentity := util.AKSPodIdentity{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
		}
		podIdentityList(AKSPodIdentity)
	},
}

var AKSPIExceptionCmd = &cobra.Command{
	Use:   "exception",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity exception`,
}

var AKSPIExceptionAddCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity exception add`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		namespace, _ := cmd.Flags().GetString("namespace")
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		podIdentityName, _ := cmd.Flags().GetString("name")
		podLabels, _ := cmd.Flags().GetString("podLabels")

		AKSPodIdentity := util.AKSPodIdentity{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Namespace:         namespace,
			PodLabels:         podLabels,
		}
		if podIdentityName != "" {
			AKSPodIdentity.Name = podIdentityName
		}
		podIdentityExceptionAdd(AKSPodIdentity)
	},
}

var AKSPIExceptionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity exception delete`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		namespace, _ := cmd.Flags().GetString("namespace")
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		podIdentityName, _ := cmd.Flags().GetString("name")

		AKSPodIdentity := util.AKSPodIdentity{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Namespace:         namespace,
			Name:              podIdentityName,
		}
		podIdentityExceptionDelete(AKSPodIdentity)
	},
}

var AKSPIExceptionListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity exception list`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")

		AKSPodIdentity := util.AKSPodIdentity{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
		}
		podIdentityExceptionList(AKSPodIdentity)
	},
}

var AKSPIExceptionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long:  `hybridctl aks pod-identity exception update`,
	Run: func(cmd *cobra.Command, args []string) {

		clusterName, _ := cmd.Flags().GetString("cluster-name")
		namespace, _ := cmd.Flags().GetString("namespace")
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		podIdentityName, _ := cmd.Flags().GetString("name")
		podLabels, _ := cmd.Flags().GetString("podLabels")

		AKSPodIdentity := util.AKSPodIdentity{
			ResourceGroupName: resourceGroupName,
			ClusterName:       clusterName,
			Namespace:         namespace,
			PodLabels:         podLabels,
			Name:              podIdentityName,
		}
		podIdentityExceptionUpdate(AKSPodIdentity)
	},
}

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long:  `hybridctl aks start --name <clusterName> --resource-group <ResourceGroupName>`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
		}
		aksStart(EKSAPIParameter)

	},
}

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long:  `hybridctl aks stop --name <clusterName> --resource-group <ResourceGroupName>`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
		}
		aksStop(EKSAPIParameter)

	},
}

var RotateCertsCmd = &cobra.Command{
	Use:   "rotate-certs",
	Short: "A brief description of your command",
	Long:  `hybridctl aks rotate-certs --name <clusterName> --resource-group <ResourceGroupName>`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
		}
		aksRotateCerts(EKSAPIParameter)
	},
}

var GetOSoptionsCmd = &cobra.Command{
	Use:   "get-os-options",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
	Run: func(cmd *cobra.Command, args []string) {

		location, _ := cmd.Flags().GetString("location")
		EKSAPIParameter := util.EKSAPIParameter{
			Location: location,
		}
		aksGetOSoptions(EKSAPIParameter)
	},
}

var MaintenanceconfigurationCmd = &cobra.Command{
	Use:   "maintenanceconfiguration",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
}

var MCAddCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
	Run: func(cmd *cobra.Command, args []string) {

		var config util.Config
		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		configName, _ := cmd.Flags().GetString("config-name")
		configFile, _ := cmd.Flags().GetString("config-file")
		// fmt.Println(configFile)
		// data, _ := ioutil.ReadFile(configFile)
		// fmt.Println(string(data))

		cobrautil.UnmarshalJsonFile(configFile, &config)

		fmt.Println(config)
		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
			ConfigName:        configName,
			ConfigFile:        config,
		}
		maintenanceconfigurationCreateOrUpdate(EKSAPIParameter)
	},
}

var MCDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		configName, _ := cmd.Flags().GetString("configname")
		if configName == "" {
			configName = "default"
		}

		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
			ConfigName:        configName,
		}
		maintenanceconfigurationDelete(EKSAPIParameter)
	},
}

var MCUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		configName, _ := cmd.Flags().GetString("configname")
		if configName == "" {
			configName = "default"
		}

		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
			ConfigName:        configName,
		}
		maintenanceconfigurationCreateOrUpdate(EKSAPIParameter)
	},
}

var MCListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		configName, _ := cmd.Flags().GetString("configname")
		if configName == "" {
			configName = "default"
		}

		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
			ConfigName:        configName,
		}
		maintenanceconfigurationList(EKSAPIParameter)
	},
}

var MCShowCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long:  `hybridctl aks get-os-options --location`,
	Run: func(cmd *cobra.Command, args []string) {

		resourceGroupName, _ := cmd.Flags().GetString("resource-group")
		clusterName, _ := cmd.Flags().GetString("name")
		configName, _ := cmd.Flags().GetString("configname")

		EKSAPIParameter := util.EKSAPIParameter{
			ResourceGroupName: resourceGroupName,
			ResourceName:      clusterName,
			ConfigName:        configName,
		}
		maintenanceconfigurationShow(EKSAPIParameter)
	},
}
