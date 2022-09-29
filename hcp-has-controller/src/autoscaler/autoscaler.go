package autoscaler

import (
	"context"
	"fmt"

	resourcev1alpha1 "hcp-pkg/apis/resource/v1alpha1"
	hasv1alpha1 "hcp-pkg/client/resource/v1alpha1/clientset/versioned"
	"hcp-pkg/util/clusterManager"

	autoscaling "k8s.io/api/autoscaling/v1"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpav1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	vpaclientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	"k8s.io/klog/v2"
)

// 각 cluster에 대한 autoscaler 저장 Map
// var AutoscalerMap map[string]*autoscaler = make(map[string]*autoscaler)
var cm, _ = clusterManager.NewClusterManager()

// cluster 단위
type Autoscaler struct {
	// clustermanager  *cm.ClusterManager
	// cluster         string
	deployList   []string
	warningcount map[string]map[string]int
	hasclientset *hasv1alpha1.Clientset
	// targetclientset *kubernetes.Clientset
}

func NewAutoScaler() *Autoscaler {
	// fmt.Printf("=> create new autoscaler for cluster %s\n", cluster)
	var hcpautoscaler Autoscaler

	// hcpautoscaler.cluster = cluster
	// hcpautoscaler.clustermanager = ncm

	config := cm.Host_config
	hasv1alpha1clientset, _ := hasv1alpha1.NewForConfig(config)
	hcpautoscaler.hasclientset = hasv1alpha1clientset
	// // hcpautoscaler.targetclientset = ncm.Cluster_kubeClients[cluster]
	// temp := make(map[*resourcev1alpha1.HCPDeployment]*resourcev1alpha1.HCPHybridAutoScaler)
	// hcpautoscaler.hasList = map[string]resourcev1alpha1.HCPHybridAutoScaler{}
	//temp2 := make(map[*resourcev1alpha1.HCPDeployment]*map[string]int)
	hcpautoscaler.warningcount = map[string]map[string]int{}
	// var list []*resourcev1alpha1.HCPDeployment
	// hcpautoscaler.deployList = list

	return &hcpautoscaler
}

// autoscaler에 해당 deployment 존재 여부 확인
func (a *Autoscaler) ExistDeployment(deployment string) bool {
	for _, dep := range a.deployList {
		if dep == deployment {
			return true
		} else {
			return false
		}
	}
	return false
}

// autoscaler에 해당 deployment 등록
func (a *Autoscaler) RegisterDeploymentToAutoScaler(deployment resourcev1alpha1.HCPDeployment) {
	fmt.Println("=> register deployment to autoscaler")

	name := deployment.ObjectMeta.Name
	a.deployList = append(a.deployList, name)
	has, _ := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, metav1.GetOptions{})
	has.Status.ScalingInProcess = true
	has.Status.ResourceStatus = "DONE"
	_, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	//a.hasList[name] = *newhas
	targets := deployment.Spec.SchedulingResult.Targets
	a.warningcount[name] = map[string]int{}
	for _, target := range targets {
		a.warningcount[name][target.Cluster] = 0
	}
}

// warningcount++
func (a *Autoscaler) WarningCountPlusOne(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target) {
	fmt.Println("=> warningcount + 1")
	name := deployment.ObjectMeta.Name
	a.warningcount[name][target.Cluster] += 1
}

func (a *Autoscaler) GetWarningCount(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target) int {
	name := deployment.ObjectMeta.Name
	return a.warningcount[name][target.Cluster]
}

/*
func (a *Autoscaler) AutoScaling(deployment *resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target) error {
	// warningcount에 따라 실행 함수 변경
	has := a.hasList[deployment]
	fmt.Println(a.hasList[deployment])
	switch a.warningcount[deployment][target.Cluster] {
	case 1:
		minReplicas := has.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas
		maxReplicas := has.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas
		return a.CreateHPA(deployment, target, minReplicas, maxReplicas)
	case 2:
		return a.UpdateHPA(deployment, target)
	case 3:
		return a.CreateVPA(deployment, target, "Auto")
	}

	return nil
}
*/
func (a *Autoscaler) CreateHPA(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target, minReplicas *int32, maxReplicas int32) error {
	fmt.Printf("===> create new HPA %s in %s\n", deployment.Name, target.Cluster)

	name := deployment.ObjectMeta.Name
	//	has := a.hasList[name]
	has, _ := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, metav1.GetOptions{})
	var namespace string
	if deployment.Spec.RealDeploymentMetadata.Namespace == "" {
		namespace = "default"
	} else {
		namespace = deployment.Spec.RealDeploymentMetadata.Namespace
	}

	// 2. hapTemplate 생성
	hpa := &hpav2beta1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Spec.RealDeploymentMetadata.Name,
			Namespace: namespace,
		},
		Spec: hpav2beta1.HorizontalPodAutoscalerSpec{
			MinReplicas: minReplicas,
			MaxReplicas: maxReplicas,
			ScaleTargetRef: hpav2beta1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       deployment.Spec.RealDeploymentMetadata.Name,
			},
		},
	}

	// 3. hpaTemplate -> HCPHybridAutoScaler 업데이트
	has.Status.LastSpec = has.Spec
	has.Spec.WarningCount = 1
	has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{HpaTemplate: *hpa}
	has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
	newhas, _ := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})

	targetclientset := cm.Cluster_kubeClients[target.Cluster]
	newhpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Create(context.TODO(), hpa, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		klog.Info("Succeed to Create HorizontalPodAutoscalers resource : ", newhpa.ObjectMeta.Name)
		newhas.Status.ResourceStatus = "DONE"
		_, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), newhas, metav1.UpdateOptions{})
		if err != nil {
			return err
		} else {
			//a.hasList[name] = has
			return nil
		}
	}
}

func (a *Autoscaler) UpdateHPA(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target) error {
	fmt.Printf("===> update HPA %s MaxReplicas in %s\n", deployment.Name, target.Cluster)
	name := deployment.ObjectMeta.Name
	//	has := a.hasList[name]
	has, _ := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, metav1.GetOptions{})
	targetclientset := cm.Cluster_kubeClients[target.Cluster]
	if has.Status.ResourceStatus == "DONE" {
		hpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(deployment.Spec.RealDeploymentMetadata.Namespace).Get(context.TODO(), deployment.Spec.RealDeploymentMetadata.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 2-1. hpa max값 설정
		nhpa := hpav2beta1.HorizontalPodAutoscaler{
			TypeMeta: metav1.TypeMeta{
				Kind:       "HorizontalPodAutoscaler",
				APIVersion: "autoscaling/v2beta1",
			},
			ObjectMeta: hpa.ObjectMeta,
			Spec: hpav2beta1.HorizontalPodAutoscalerSpec{
				MinReplicas:    hpa.Spec.MinReplicas,
				MaxReplicas:    hpa.Spec.MaxReplicas * 2,
				ScaleTargetRef: hpa.Spec.ScaleTargetRef,
			},
		}
		// 3. hpaTemplate -> HCPHybridAutoScaler 업데이트
		has.Status.LastSpec = has.Spec
		has.Spec.WarningCount = 2
		has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{HpaTemplate: nhpa}
		has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
		//a.hasList[name] = has
		newhas, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}

		targetclientset := cm.Cluster_kubeClients[target.Cluster]
		newhpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(deployment.Spec.RealDeploymentMetadata.Namespace).Update(context.TODO(), hpa, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			klog.Info("Succeed to Create HorizontalPodAutoscalers resource : ", newhpa.ObjectMeta.Name)
			newhas.Status.ResourceStatus = "DONE"
			_, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), newhas, metav1.UpdateOptions{})
			if err != nil {
				return err
			} else {
				//a.hasList[name] = has
				fmt.Printf("=====> update %s Done\n", newhas.Name)
				return nil
			}
		}
	} else {
		fmt.Printf("HCPHybridAutoScaler ResourceStatus is not DONE : %s\n", has.Status.ResourceStatus)
		return fmt.Errorf("HCPHybridAutoScaler ResourceStatus is not DONE : %s", has.Status.ResourceStatus)
	}
}

func (a *Autoscaler) CreateVPA(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target, updateMode string) error {
	fmt.Printf("===> create new VPA %s in %s\n", deployment.Name, target.Cluster)
	name := deployment.ObjectMeta.Name
	has, _ := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, metav1.GetOptions{})
	//has := a.hasList[name]
	if has.Status.ResourceStatus == "DONE" {
		// 2. vpaTemplate 생성
		// updateMode := vpav1beta2.UpdateModeAuto
		vpa := vpav1beta2.VerticalPodAutoscaler{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "autoscaling.k8s.io/v1",
				Kind:       "VerticalPodAutoscaler",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployment.Name,
				Namespace: deployment.Namespace,
			},
			Spec: vpav1beta2.VerticalPodAutoscalerSpec{
				TargetRef: &autoscaling.CrossVersionObjectReference{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       deployment.Name,
				},
				UpdatePolicy: &vpav1beta2.PodUpdatePolicy{
					UpdateMode: (*vpav1beta2.UpdateMode)(&updateMode),
				},
			},
		}

		has, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 3. vpaTemplate -> HCPHybridAutoScaler 생성
		has.Status.LastSpec = has.Spec
		has.Spec.WarningCount = 3
		has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{VpaTemplate: vpa}
		has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
		//	a.hasList[name] = *has
		newhas, _ := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})

		target_config := cm.Cluster_configs[target.Cluster]
		vpa_clientset, _ := vpaclientset.NewForConfig(target_config)
		hcphas, err := vpa_clientset.AutoscalingV1beta2().VerticalPodAutoscalers(deployment.Spec.RealDeploymentMetadata.Namespace).Create(context.TODO(), &vpa, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			klog.Info("Success to Create VerticalPodAutoscalers resource : ", hcphas.ObjectMeta.Name)
			newhas.Status.ResourceStatus = "DONE"
			_, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), newhas, metav1.UpdateOptions{})
			if err != nil {
				return err
			} else {
				//a.hasList[name] = *has
				fmt.Printf("=====> update %s Done\n", newhas.Name)
				return nil
			}
		}
	} else {
		fmt.Printf("HCPHybridAutoScaler ResourceStatus is not DONE : %s\n", has.Status.ResourceStatus)
		return fmt.Errorf("HCPHybridAutoScaler ResourceStatus is not DONE : %s", has.Status.ResourceStatus)
	}
}
