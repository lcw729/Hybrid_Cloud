package cmd

func aksFlags() {
	aksCmd.PersistentFlags().StringP("resource-group", "g", "", "Name of resource group.")
	aksCmd.PersistentFlags().StringP("name", "n", "", "Name of the managed cluster.")

	//addon
	AddonCmd.PersistentFlags().StringP("addon", "a", "", "Specify the Kubernetes addon")

	AKSDisableAddonsCmd.MarkFlagRequired("resource-group")
	AKSDisableAddonsCmd.MarkFlagRequired("name")
	AKSDisableAddonsCmd.MarkFlagRequired("addon")

	AKSEnableAddonsCmd.MarkFlagRequired("resource-group")
	AKSEnableAddonsCmd.MarkFlagRequired("name")
	AKSEnableAddonsCmd.MarkFlagRequired("addon")

	AKSListAddonsCmd.MarkFlagRequired("resource-group")
	AKSListAddonsCmd.MarkFlagRequired("name")

	AKSShowAddonsCmd.MarkFlagRequired("resource-group")
	AKSShowAddonsCmd.MarkFlagRequired("name")
	AKSShowAddonsCmd.MarkFlagRequired("addon")

	AKSUpdateAddonsCmd.MarkFlagRequired("resource-group")
	AKSUpdateAddonsCmd.MarkFlagRequired("name")
	AKSUpdateAddonsCmd.MarkFlagRequired("addon")

	//Pod-Identity

	AKSPodIdentityCmd.PersistentFlags().String("cluster-name", "", "The cluster name.")
	AKSPodIdentityCmd.PersistentFlags().String("namespace", "", "The pod identity namespace.")
	AKSPodIdentityCmd.PersistentFlags().StringP("name", "", "n", "The pod identity name. Generate if not specified.")

	AKSPIAddCmd.MarkPersistentFlagRequired("resource-group")
	AKSPIAddCmd.MarkPersistentFlagRequired("cluster-name")
	AKSPIAddCmd.MarkPersistentFlagRequired("namespace")
	AKSPIAddCmd.MarkPersistentFlagRequired("name")
	AKSPIAddCmd.Flags().String("identity-resource-id", "", "Resource id of the identity to use.")
	AKSPIAddCmd.Flags().String("binding-selector", "", "Optional binding selector to use.")

	AKSPIDeleteCmd.MarkPersistentFlagRequired("cluster-name")
	AKSPIDeleteCmd.MarkPersistentFlagRequired("namespace")
	AKSPIDeleteCmd.MarkPersistentFlagRequired("name")
	AKSPIDeleteCmd.MarkPersistentFlagRequired("resource-group")

	AKSPIListCmd.MarkPersistentFlagRequired("cluster-name")
	AKSPIListCmd.MarkPersistentFlagRequired("resource-group")

	AKSPIExceptionCmd.Flags().String("pod-labels", "", "Space-separated labels: key=value [key=value ...].")
	AKSPIExceptionAddCmd.MarkPersistentFlagRequired("cluster-name")
	AKSPIExceptionAddCmd.MarkPersistentFlagRequired("namespace")
	AKSPIExceptionAddCmd.MarkPersistentFlagRequired("pod-labels")
	AKSPIExceptionAddCmd.MarkPersistentFlagRequired("resource-group")
	AKSPIExceptionAddCmd.MarkPersistentFlagRequired("name")

	AKSPIExceptionDeleteCmd.MarkPersistentFlagRequired("cluster-name")
	AKSPIExceptionDeleteCmd.MarkPersistentFlagRequired("name")
	AKSPIExceptionDeleteCmd.MarkPersistentFlagRequired("namespace")
	AKSPIExceptionDeleteCmd.MarkPersistentFlagRequired("resource-group")

	AKSPIExceptionListCmd.MarkPersistentFlagRequired("cluster-name")
	AKSPIExceptionListCmd.MarkPersistentFlagRequired("resource-group")

	AKSPIExceptionUpdateCmd.MarkPersistentFlagRequired("cluster-name")
	AKSPIExceptionUpdateCmd.MarkPersistentFlagRequired("name")
	AKSPIExceptionUpdateCmd.MarkPersistentFlagRequired("namespace")
	AKSPIExceptionUpdateCmd.MarkPersistentFlagRequired("resource-group")
	AKSPIExceptionUpdateCmd.MarkPersistentFlagRequired("pod-labels")

	StopCmd.MarkPersistentFlagRequired("resource-group")
	StopCmd.MarkPersistentFlagRequired("name")

	GetOSoptionsCmd.PersistentFlags().StringP("location", "l", "", "location")
	GetOSoptionsCmd.MarkPersistentFlagRequired("location")

	MaintenanceconfigurationCmd.MarkPersistentFlagRequired("resource-group")
	MaintenanceconfigurationCmd.MarkPersistentFlagRequired("name")

	MCAddCmd.Flags().StringP("config-name", "c", "", "configname")
	MCAddCmd.MarkFlagRequired("config-name")
	MCAddCmd.Flags().StringP("config-file", "", "", "configfile")
	MCAddCmd.MarkFlagRequired("config-file")

	MCDeleteCmd.Flags().StringP("config-name", "c", "", "configname")
	MCDeleteCmd.MarkFlagRequired("config-name")

	MCUpdateCmd.Flags().StringP("config-name", "c", "", "configname")
	MCUpdateCmd.MarkFlagRequired("config-name")
	MCUpdateCmd.Flags().StringP("config-file", "", "", "configfile")
	MCUpdateCmd.MarkFlagRequired("config-file")

	MCShowCmd.Flags().StringP("config-name", "c", "", "configname")
	MCShowCmd.MarkFlagRequired("config-name")

}
