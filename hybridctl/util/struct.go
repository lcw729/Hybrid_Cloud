package util

type EksAPIParameter struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceName      string
	ApiVersion        string
}

type Addon struct {
	AddonName     string
	ClusterName   string
	Message_      string
	NodegroupName string
}
