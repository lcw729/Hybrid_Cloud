package Policy

import (
	"Hybrid_Cluster/clientset/v1alpha1"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cobrautil "Hybrid_Cluster/hybridctl/util"
	"strconv"
)

func GetCycle() float64 {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := v1alpha1.NewForConfig(master_config)
	hcppolicy, err := clientset.HCPPolicy("hcp").Get("weight-calculation-cycle", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	policies := hcppolicy.Spec.Template.Spec.Policies
	if len(policies) == 1 {
		policy := policies[0]
		if policy.Type == "cycle" {
			cycle, _ := strconv.ParseFloat(policy.Value[0], 64)
			if cycle > 0 {
				return cycle
			}
		}
	}
	return -1
}
