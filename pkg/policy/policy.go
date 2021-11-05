package policy

import (
	"Hybrid_Cluster/clientset/v1alpha1"
	"fmt"
	"log"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	util "Hybrid_Cluster/hcp-scheduler/pkg/util"
	cobrautil "Hybrid_Cluster/hybridctl/util"
)

func GetWatchingLevel() util.WatchingLevel {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := v1alpha1.NewForConfig(master_config)
	hcppolicy, err := clientset.HCPPolicy("hcp").Get("watching-level", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	var watchingLevel util.WatchingLevel
	watchingLevel.Levels = make([]util.Level, 5)
	for i, policy := range hcppolicy.Spec.Template.Spec.Policies {
		var level util.Level
		level = util.Level{
			Type:  policy.Type,
			Value: policy.Value,
		}
		watchingLevel.Levels[i] = level
	}
	return watchingLevel
}

func GetWarningLevel() util.Level {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := v1alpha1.NewForConfig(master_config)
	hcppolicy, err := clientset.HCPPolicy("hcp").Get("warning-level", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	// for i, policy := range hcppolicy.Spec.Template.Spec.Policies {
	policies := hcppolicy.Spec.Template.Spec.Policies
	if len(policies) == 1 {
		level := util.Level{
			Type:  hcppolicy.Spec.Template.Spec.Policies[0].Type,
			Value: hcppolicy.Spec.Template.Spec.Policies[0].Value,
		}
		return level
	} else {
		level := util.Level{
			Type: "none",
		}
		return level
	}
}

func GetCycle() float64 {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := v1alpha1.NewForConfig(master_config)
	hcppolicy, err := clientset.HCPPolicy("hcp").Get("weight-calculation-cycle", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	for _, policy := range hcppolicy.Spec.Template.Spec.Policies {
		if policy.Type == "cycle" && len(policy.Value) == 1 {
			cycle, err := strconv.ParseFloat(policy.Value[0], 64)
			if err == nil && cycle > 0 {
				fmt.Println("Policy Type : ", "cycle")
				fmt.Println("Policy Value [cycle] : ", cycle)
				return cycle
			}
		}
	}
	return -1
}
