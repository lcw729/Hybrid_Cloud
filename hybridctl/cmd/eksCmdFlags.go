package cmd

func eksFlags() {
	EKSAssociateEncryptionConfigCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster that you are associating with encryption configuration.")
	EKSAssociateEncryptionConfigCmd.MarkFlagRequired("cluster-name")
	EKSAssociateEncryptionConfigCmd.Flags().StringP("encryption-config", "", "", "The configuration you are using for encryption.")
	EKSAssociateEncryptionConfigCmd.MarkFlagRequired("encryption-config")
	EKSAssociateEncryptionConfigCmd.Flags().StringP("client-request-token", "", "", "The client request token you are using with the encryption configuration.")

	EKSAssociateIdentityProviderConfigCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster to associate the configuration to.")
	EKSAssociateIdentityProviderConfigCmd.MarkFlagRequired("cluster-name")
	EKSAssociateIdentityProviderConfigCmd.Flags().StringP("oidc", "", "", "An object that represents an OpenID Connect (OIDC) identity provider configuration.")
	EKSAssociateIdentityProviderConfigCmd.MarkFlagRequired("oidc")
	EKSAssociateIdentityProviderConfigCmd.Flags().StringP("client-request-token", "", "", "enter client request token")
	EKSAssociateIdentityProviderConfigCmd.Flags().StringP("tags", "", "", "enter your tags Jsonfile name")

	EKSCreateAddonCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster to create the add-on for.")
	EKSCreateAddonCmd.MarkFlagRequired("cluster-name")
	EKSCreateAddonCmd.Flags().StringP("addon-name", "a", "", "The name of the add-on. The name must match one of the names returned by DescribeAddonVersions")
	EKSCreateAddonCmd.MarkFlagRequired("addon-name")
	EKSCreateAddonCmd.Flags().StringP("addon-version", "", "", "The version of the add-on. The version must match one of the versions returned by DescribeAddonVersions")
	EKSCreateAddonCmd.Flags().StringP("service-account-role-arn", "", "", "The Amazon Resource Name (ARN) of an existing IAM role to bind to the add-on's service account.")
	EKSCreateAddonCmd.Flags().StringP("resolve-conflicts", "", "", "How to resolve parameter value conflicts when migrating an existing add-on to an Amazon EKS add-on. Possible values: OVERWRITE, NONE")
	EKSCreateAddonCmd.Flags().StringP("client-request-token", "", "", "A unique, case-sensitive identifier that you provide to ensure the idempotency of the request.")
	EKSCreateAddonCmd.Flags().StringP("tags", "", "", "The metadata to apply to the cluster to assist with categorization and organization. Shorthand Syntax: KeyName1=string,KeyName2=string")

	EKSDeleteAddonCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster to delete the add-on from.")
	EKSDeleteAddonCmd.MarkFlagRequired("cluster-name")
	EKSDeleteAddonCmd.Flags().StringP("addon-name", "a", "", "The name of the add-on. The name must match one of the names returned by ListAddons")
	EKSDeleteAddonCmd.MarkFlagRequired("addon-name")

	EKSDescribeAddonCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster.")
	EKSDescribeAddonCmd.MarkFlagRequired("cluster-name")
	EKSDescribeAddonCmd.Flags().StringP("addon-name", "a", "", "The name of the add-on. The name must match one of the names returned by ListAddons")
	EKSDescribeAddonCmd.MarkFlagRequired("addon-name")

	EKSDescribeAddonVersionsCmd.Flags().StringP("addon-name", "a", "", "The name of the add-on. The name must match one of the names returned by ListAddons")
	EKSDescribeAddonVersionsCmd.Flags().StringP("kubernetes-version", "", "", "The Kubernetes versions that the add-on can be used with.")
	EKSDescribeAddonVersionsCmd.Flags().Int64P("max-results", "", 0, "The maximum number of results to return.")
	EKSDescribeAddonVersionsCmd.Flags().StringP("next-token", "", "", "The nextToken value returned from a previous paginated DescribeAddonVersionsRequest where maxResults was used and the results exceeded the value of that parameter.")

	describeUpdateCmd.Flags().StringP("name", "c", "", "The name of the Amazon EKS cluster associated with the update.")
	describeUpdateCmd.MarkFlagRequired("name")
	describeUpdateCmd.Flags().StringP("update-id", "", "", "")
	describeUpdateCmd.MarkFlagRequired("update-id")
	describeUpdateCmd.Flags().StringP("nodegroup-name", "", "", "enter nodegroupName")
	describeUpdateCmd.Flags().StringP("addon-name", "", "", "enter addonName")

	EKSDisassociateIdentityProviderConfigCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster to disassociate an identity provider from.")
	EKSDisassociateIdentityProviderConfigCmd.MarkFlagRequired("cluster-name")
	EKSDisassociateIdentityProviderConfigCmd.Flags().StringP("identity-provider-config", "", "", "An object that represents an identity provider configuration.")
	EKSDisassociateIdentityProviderConfigCmd.MarkFlagRequired("identity-provider-config")
	EKSDisassociateIdentityProviderConfigCmd.Flags().StringP("client-request-token", "", "", "A unique, case-sensitive identifier that you provide to ensure the idempotency of the request.")

	EKSListAddonCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster.")
	EKSListAddonCmd.MarkFlagRequired("cluster-name")
	EKSListAddonCmd.Flags().Int64P("max-result", "", 0, "The maximum number of add-on results returned by ListAddonsRequest in paginated")
	EKSListAddonCmd.Flags().StringP("next-token", "", "", "The nextToken value returned from a previous paginated ListAddonsRequest")

	EKSListIdentityProviderConfigsCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster.")
	EKSListIdentityProviderConfigsCmd.MarkFlagRequired("cluster-name")
	EKSListIdentityProviderConfigsCmd.Flags().Int64P("max-result", "", 0, "enter maxresult")
	EKSListIdentityProviderConfigsCmd.Flags().StringP("next-token", "", "", "enter next token")

	EKSDescribeIdentityProviderConfigCmd.Flags().StringP("cluster-name", "c", "", "The cluster name that the identity provider configuration is associated to.")
	EKSDescribeIdentityProviderConfigCmd.MarkFlagRequired("cluster-name")
	EKSDescribeIdentityProviderConfigCmd.Flags().StringP("identity-provider-config", "", "", "An object that represents an identity provider configuration.")
	EKSDescribeIdentityProviderConfigCmd.MarkFlagRequired("identity-provider-config")

	EKSListTagsForResourceCmd.Flags().StringP("resource-arn", "", "", "Enter resource-arn")

	listUpdateCmd.Flags().StringP("name", "c", "", "The name of the Amazon EKS cluster associated with the update.")
	listUpdateCmd.MarkFlagRequired("name")
	listUpdateCmd.Flags().StringP("nodegroup-name", "", "", "enter nodegroupName")
	listUpdateCmd.Flags().StringP("addon-name", "", "", "enter addonName")
	listUpdateCmd.Flags().Int64P("max-result", "", 0, "enter maxresult")
	listUpdateCmd.Flags().StringP("next-token", "", "", "enter next token")

	EKSTagResourceCmd.Flags().StringP("tags", "t", "", "enter your tags Jsonfile name")
	EKSTagResourceCmd.MarkPersistentFlagRequired("tags")
	EKSTagResourceCmd.Flags().StringP("resource-arn", "", "", "Enter resource-arn")
	EKSTagResourceCmd.MarkPersistentFlagRequired("resource-arn")

	EKSUntagResourceCmd.Flags().StringP("resource-arn", "", "", "Enter resource-arn")
	EKSUntagResourceCmd.Flags().StringP("tag-keys", "t", "", "enter your tag-keys list")
	EKSUntagResourceCmd.MarkPersistentFlagRequired("tag-keys")
	EKSUntagResourceCmd.MarkPersistentFlagRequired("resource-arn")

	EKSUpdateAddonCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster.")
	EKSUpdateAddonCmd.MarkFlagRequired("cluster-name")
	EKSUpdateAddonCmd.Flags().StringP("addon-name", "a", "", "The name of the add-on. The name must match one of the names returned by ListAddons")
	EKSUpdateAddonCmd.MarkFlagRequired("addon-name")
	EKSUpdateAddonCmd.Flags().StringP("addon-version", "", "", "The version of the add-on. The version must match one of the versions returned by DescribeAddonVersions")
	EKSUpdateAddonCmd.Flags().StringP("service-account-role-arn", "", "", "The Amazon Resource Name (ARN) of an existing IAM role to bind to the add-on's service account.")
	EKSUpdateAddonCmd.Flags().StringP("resolve-conflicts", "", "", "How to resolve parameter value conflicts when migrating an existing add-on to an Amazon EKS add-on. Possible values: OVERWRITE, NONE")
	EKSUpdateAddonCmd.Flags().StringP("client-request-token", "", "", "Unique, case-sensitive identifier that you provide to ensure the idempotency of the request.")

	EKSUpdateClusterConfigCmd.Flags().StringP("name", "c", "", "The name of the Amazon EKS cluster associated with the update.")
	EKSUpdateClusterConfigCmd.MarkFlagRequired("name")
	EKSUpdateClusterConfigCmd.Flags().StringP("resource-vpc-config", "", "", "enter resource-vpc-config jsonfile name")
	EKSUpdateClusterConfigCmd.Flags().StringP("logging", "", "", "enter logging jsonfile name")
	EKSUpdateClusterConfigCmd.Flags().StringP("client-request-token", "", "", "enter client request token")

	EKSUpdateNodegroupConfigCmd.Flags().StringP("cluster-name", "c", "", "The name of the cluster.")
	EKSUpdateNodegroupConfigCmd.MarkFlagRequired("cluster-name")
	EKSUpdateNodegroupConfigCmd.Flags().StringP("nodegroup-name", "", "", "enter nodegroupName")
	EKSUpdateNodegroupConfigCmd.MarkFlagRequired("nodegroup-name")
	EKSUpdateNodegroupConfigCmd.Flags().StringP("labels", "", "", "enter labels jsonfile name")
	EKSUpdateNodegroupConfigCmd.Flags().StringP("taints", "", "", "enter taints jsonfile name")
	EKSUpdateNodegroupConfigCmd.Flags().StringP("scaling-config", "", "", "enter resource-vpc-config jsonfile name")
	EKSUpdateNodegroupConfigCmd.Flags().StringP("update-config", "", "", "enter logging jsonfile name")
	EKSUpdateNodegroupConfigCmd.Flags().StringP("client-request-token", "", "", "enter client request token")

}
