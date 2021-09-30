package cmd

func aksFlags() {
	aksCmd.PersistentFlags().StringP("resource-group", "g", "", "resourceGroup name")
	aksCmd.PersistentFlags().StringP("name", "n", "", "clustername")
	aksCmd.MarkPersistentFlagRequired("resource-group")
	aksCmd.MarkPersistentFlagRequired("name")
}
