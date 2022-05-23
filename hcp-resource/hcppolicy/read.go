package resource

import (
	"Hybrid_Cloud/hcp-analytic-engine/util"
	hcppolicyapis "Hybrid_Cloud/pkg/apis/hcppolicy/v1alpha1"
	hcppolicyv1alpha1 "Hybrid_Cloud/pkg/client/hcppolicy/v1alpha1/clientset/versioned"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type WatchingLevel struct {
	Levels []Level
}

type Level struct {
	Type  string
	Value string
}

func GetHCPPolicy(policy_name string) (*hcppolicyapis.HCPPolicy, error) {
	cm, err := clusterManager.NewClusterManager()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	c, err := hcppolicyv1alpha1.NewForConfig(cm.Host_config)
	if err != nil {
		klog.Info(err)
		return nil, err
	}
	policy, err := c.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), policy_name, metav1.GetOptions{})
	if err != nil {
		klog.Info(err)
		return nil, err
	}
	return policy, nil
}

func GetInitialSettingValue(typ string) (int, string) {
	policy, err := GetHCPPolicy("initial-setting")
	if err != nil {
		klog.Info(err)
	}

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

func GetWatchingLevel() WatchingLevel {
	hcppolicy, err := GetHCPPolicy("watching-level")
	if err != nil {
		klog.Info(err)
	}

	var watchingLevel WatchingLevel
	watchingLevel.Levels = make([]Level, 5)
	for i, policy := range hcppolicy.Spec.Template.Spec.Policies {
		var level Level = Level{
			Type:  policy.Type,
			Value: policy.Value,
		}
		watchingLevel.Levels[i] = level
	}
	return watchingLevel
}

func GetWarningLevel() Level {
	hcppolicy, err := GetHCPPolicy("warning-level")
	if err != nil {
		klog.Info(err)
	}

	// for i, policy := range hcppolicy.Spec.Template.Spec.Policies {
	policies := hcppolicy.Spec.Template.Spec.Policies
	if len(policies) == 1 {
		level := Level{
			Type:  hcppolicy.Spec.Template.Spec.Policies[0].Type,
			Value: hcppolicy.Spec.Template.Spec.Policies[0].Value,
		}
		return level
	} else {
		level := Level{
			Type: "none",
		}
		return level
	}
}

func GetAlgorithm() (string, error) {
	hcppolicy, err := GetHCPPolicy("optimal-arrangement-algorithm")
	if err != nil {
		klog.Info(err)
	}

	policies := hcppolicy.Spec.Template.Spec.Policies
	if len(policies) == 1 {
		policy := policies[0]
		if policy.Type == "Algorithm" {
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
	hcppolicy, err := GetHCPPolicy("weight-calculation-cycle")
	if err != nil {
		klog.Info(err)
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
