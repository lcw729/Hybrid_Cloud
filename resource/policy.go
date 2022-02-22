package resource

import (
	"Hybrid_Cluster/hcp-analytic-engine/util"
	hcppolicyapis "Hybrid_Cluster/pkg/apis/hcppolicy/v1alpha1"
	hcppolicyv1alpha1 "Hybrid_Cluster/pkg/client/hcppolicy/v1alpha1/clientset/versioned"
	"Hybrid_Cluster/util/clusterManager"
	"context"
	"fmt"
	"log"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	cobrautil "Hybrid_Cluster/hybridctl/util"
)

func GetInitialSettingValue(typ string) (int, string) {
	policy := GetPolicy("initial-setting")
	policies := policy.Spec.Template.Spec.Policies
	for _, p := range policies {
		println(p.Type)
		if typ == "default_node_option" && p.Type == "default_node_option" {
			var value string = p.Value
			if value == "" {
				fmt.Printf("ERROR: No %s Value\n", typ)
			} else {
				return -1, value
			}
		} else if p.Type == typ {
			var value int
			value, err := strconv.Atoi(p.Value)
			if err != nil {
				fmt.Printf("ERROR: No %s Value\n", typ)
			}
			return value, ""
		}
	}
	fmt.Printf("ERROR: No Such Type %s\n", typ)
	return -1, ""
}

func GetPolicy(policy_name string) *hcppolicyapis.HCPPolicy {
	cm := clusterManager.NewClusterManager()

	c, err := hcppolicyv1alpha1.NewForConfig(cm.Host_config)
	if err != nil {
		klog.Info(err)
	}
	policy, err := c.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), policy_name, metav1.GetOptions{})
	if err != nil {
		klog.Info(err)
	}
	return policy
}

func GetWatchingLevel() util.WatchingLevel {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := hcppolicyv1alpha1.NewForConfig(master_config)
	if err != nil {
		fmt.Println(err)
	}
	hcppolicy, err := clientset.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), "watching-level", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}

	var watchingLevel util.WatchingLevel
	watchingLevel.Levels = make([]util.Level, 5)
	for i, policy := range hcppolicy.Spec.Template.Spec.Policies {
		var level util.Level = util.Level{
			Type:  policy.Type,
			Value: policy.Value,
		}
		watchingLevel.Levels[i] = level
	}
	return watchingLevel
}

func GetWarningLevel() util.Level {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := hcppolicyv1alpha1.NewForConfig(master_config)
	if err != nil {
		log.Println(err)
	}
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

func GetAlgorithm() (string, error) {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := hcppolicyv1alpha1.NewForConfig(master_config)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	hcppolicy, err := clientset.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), "optimal-arrangement-algorithm", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	policies := hcppolicy.Spec.Template.Spec.Policies
	if len(policies) == 1 {
		policy := policies[0]
		if policy.Type == "Algorithm" && len(policy.Value) == 1 {
			algorithm := policy.Value
			for _, algo := range util.AlgorithmList {
				if algo == algorithm {
					return algorithm, nil
				}
			}
		}
	}
	return "", err
}

func GetCycle() float64 {
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, err := hcppolicyv1alpha1.NewForConfig(master_config)
	if err != nil {
		log.Println(err)
	}
	hcppolicy, err := clientset.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), "weight-calculation-cycle", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	for _, policy := range hcppolicy.Spec.Template.Spec.Policies {
		if policy.Type == "cycle" && len(policy.Value) == 1 {
			cycle, err := strconv.ParseFloat(policy.Value, 64)
			if err == nil && cycle > 0 {
				fmt.Println("Policy Type : ", "cycle")
				fmt.Println("Policy Value [cycle] : ", cycle)
				return cycle
			}
		}
	}
	return -1
}
