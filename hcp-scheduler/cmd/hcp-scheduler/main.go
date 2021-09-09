package main

import (
	// "Hybrid_Cluster/analytic-engine/analyticEngine"

	analyticEngine "Hybrid_Cluster/hcp-analytic-engine/analyticEngine"
	scheduler "Hybrid_Cluster/hcp-scheduler/pkg"
	"fmt"
	"strconv"
	"time"
)

func main() {

	// host_ctx := "kube-master"
	// namespace := "hybrid"
	// host_cfg, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	// // host_cfg := cm.Host_config
	// live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})

	// ghosts := []*cluster.Cluster{}

	// co, _ := templateresource.NewController(live, ghosts, namespace)

	// m := manager.New()
	// m.AddController(co)
	// if err := m.Start(signals.SetupSignalHandler()); err != nil {
	// 	log.Fatal(err)
	// }
	// ResourceExtensionSchedule()
	ResourceExtensionSchedule()
	// time := time.Second * 2

}

func ResourceExtensionSchedule() {
	fmt.Println("[ Scheduler Start ]")
	fmt.Println("[step 1] Check Policy from Policy manager - calculation cycle")
	period := getPeriod()
	fmt.Println("")
	if period > 0 {
		fmt.Println("-------------------------LOOP START----------------------------")
		for {
			time.Sleep(time.Second * time.Duration(period))
			fmt.Println("Resource Extension Call")
			analyticEngine.ResourceExtension()
			scheduler.Resourcebalancingcontroller()
			fmt.Println("---------------------------------------------------------------")
		}
	}
}

func optimalArrangement() {
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[ Scheduler Start ]")
	fmt.Println("[step 1] Check Policy from Policy manager - DRF, Affinity")
	fmt.Println("----> Policy Value: Affinity")
	fmt.Println("----> Policy Value: DRF")
	fmt.Println("[step 2] Call Affinity Calulator")
	analyticEngine.AffinityCalculator()
	fmt.Println("[step 3] Profiling Pod & Node")
	analyticEngine.DRF()
	fmt.Println("[step 4] Checking Pending POD Queue")
	fmt.Println("----> [case 1] If there are no suitable resources, wait ")
	fmt.Println("----> [case 2] Select suitable target resources and resource request ")
}

func getPeriod() int {
	// master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	// clientset, err := v1alpha1.NewForConfig(master_config)
	// hcppolicy, err := clientset.HCPPolicy("hybrid").Get("resource-weight-period", metav1.GetOptions{})
	// if err != nil {
	// 	log.Println(err)
	// }
	// for _, policy := range hcppolicy.Spec.Template.Spec.Policies {
	// 	if policy.Type == "period" && len(policy.Value) == 1 {
	period, err := strconv.Atoi("2")
	if err == nil && period > 0 {
		fmt.Println("Policy Type : ", "Period")
		fmt.Println("Policy Value [Period] : ", period)
		return period
	}
	// 	}
	// }
	return -1
}
