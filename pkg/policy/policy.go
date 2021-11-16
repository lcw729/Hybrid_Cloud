package policy

import (
	"Hybrid_Cluster/clientset/v1alpha1"
	"fmt"
	"log"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cobrautil "Hybrid_Cluster/hybridctl/util"
)

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
