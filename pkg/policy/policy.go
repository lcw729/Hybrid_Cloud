package Policy

import (
	"Hybrid_Cluster/hcp-analytic-engine/util"
	"context"
	"log"

	v1alpha1 "Hybrid_Cluster/pkg/client/policy/v1alpha1/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cobrautil "Hybrid_Cluster/hybridctl/util"
)

// func GetCycle() float64 {
// 	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
// 	clientset, err := v1alpha1.NewForConfig(master_config)
// 	hcppolicy, err := clientset.HCPPolicy("hcp").Get("weight-calculation-cycle", metav1.GetOptions{})
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	policies := hcppolicy.Spec.Template.Spec.Policies
// 	if len(policies) == 1 {
// 		policy := policies[0]
// 		if policy.Type == "cycle" {
// 			cycle, _ := strconv.ParseFloat(policy.Value[0], 64)
// 			if cycle > 0 {
// 				return cycle
// 			}
// 		}
// 	}
// 	return -1
// }

func GetWatchingLevel() util.WatchingLevel {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := v1alpha1.NewForConfig(master_config)
	hcppolicy, err := clientset.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), "watching-level", metav1.GetOptions{})
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
	hcppolicy, err := clientset.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), "warning-level", metav1.GetOptions{})
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

func GetAlgorithm() string {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := v1alpha1.NewForConfig(master_config)
	hcppolicy, err := clientset.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), "optimal-arrangement-algorithm", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	policies := hcppolicy.Spec.Template.Spec.Policies
	if len(policies) == 1 {
		policy := policies[0]
		if policy.Type == "Algorithm" && len(policy.Value) == 1 {
			algorithm := policy.Value[0]
			for _, algo := range util.AlgorithmList {
				if algo == algorithm {
					return algorithm
				}
			}
		}
	}
	return ""
}
