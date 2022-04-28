package resource

import (
	resourcev1alpha1 "Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	hasv1alpha1 "Hybrid_Cloud/pkg/client/resource/v1alpha1/clientset/versioned"
	cm "Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpav1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	"k8s.io/client-go/kubernetes"
)

var AutoscalerMap map[string]*autoscaler = make(map[string]*autoscaler)

// cluster 단위
type autoscaler struct {
	clustermanager *cm.ClusterManager
	cluster        string
	//deployment      *appsv1.Deployment
	has_name        *map[*appsv1.Deployment]string
	warningcount    *map[*appsv1.Deployment]int
	hasclientset    *hasv1alpha1.Clientset
	targetclientset *kubernetes.Clientset
}

func NewAutoScaler(cluster string, deployment *appsv1.Deployment, namespace string) *autoscaler {
	var hcpautoscaler autoscaler

	hcpautoscaler.cluster = cluster

	ncm, _ := cm.NewClusterManager()
	hcpautoscaler.clustermanager = ncm

	config := ncm.Host_config
	hasv1alpha1clientset, _ := hasv1alpha1.NewForConfig(config)
	hcpautoscaler.hasclientset = hasv1alpha1clientset
	hcpautoscaler.targetclientset = ncm.Cluster_kubeClients[cluster]

	//hcpautoscaler.deployment = deployment
	name := make(map[*appsv1.Deployment]string)
	hcpautoscaler.has_name = &name
	(*hcpautoscaler.has_name)[deployment] = cluster + "-" + deployment.Name

	has, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), (*hcpautoscaler.has_name)[deployment], metav1.GetOptions{})
	if err != nil {
		temp := make(map[*appsv1.Deployment]int)
		hcpautoscaler.warningcount = &temp
		(*hcpautoscaler.warningcount)[deployment] = 0
	} else {
		temp := make(map[*appsv1.Deployment]int)
		hcpautoscaler.warningcount = &temp
		(*hcpautoscaler.warningcount)[deployment] = int(has.Spec.WarningCount)
	}

	return &hcpautoscaler
}

func (a *autoscaler) WarningCountPlusOne(deployment *appsv1.Deployment) {
	(*a.warningcount)[deployment] += 1
}

func (a *autoscaler) GetWarningCount(deployment *appsv1.Deployment) int {
	return (*a.warningcount)[deployment]
}

func (a *autoscaler) AutoScaling(deployment *appsv1.Deployment) error {
	// 2. HCPHybridAutoScalers 정보 얻기
	switch (*a.warningcount)[deployment] {
	case 1:
		var min int32 = 1
		minReplicas := &min
		var maxReplicas int32 = 5
		return a.CreateHPA(deployment, minReplicas, maxReplicas)
	case 2:
		return a.UpdateHPA(deployment)
	case 3:
		return a.CreateVPA(deployment, "Auto")
	}

	return nil
}

func (a *autoscaler) CreateHPA(deployment *appsv1.Deployment, minReplicas *int32, maxReplicas int32) error {

	// 2. hapTemplate 생성
	hpa := &hpav2beta1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		},
		Spec: hpav2beta1.HorizontalPodAutoscalerSpec{
			MinReplicas: minReplicas,
			MaxReplicas: maxReplicas,
			ScaleTargetRef: hpav2beta1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       deployment.Name,
			},
		},
	}

	// 3. hpaTemplate -> HCPHybridAutoScaler 생성
	instance := &resourcev1alpha1.HCPHybridAutoScaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: a.cluster + "-" + hpa.Spec.ScaleTargetRef.Name,
		},
		Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
			TargetCluster: a.cluster,
			WarningCount:  1,
			ScalingOptions: resourcev1alpha1.ScalingOptions{
				HpaTemplate: *hpa,
			},
		},
		Status: resourcev1alpha1.HCPHybridAutoScalerStatus{
			ResourceStatus: "WAITING",
		},
	}

	newhas, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Create(context.TODO(), instance, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("create %s Done\n", newhas.Name)
		return nil
	}
}

func (a *autoscaler) UpdateHPA(deployment *appsv1.Deployment) error {
	hpa, err := a.targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(deployment.Namespace).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
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

	// has_name := a.cluster + "-" + a.deployment.Name
	has, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), (*a.has_name)[deployment], metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	has.Status.LastSpec = has.Spec
	has.Spec.WarningCount = 2
	has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{HpaTemplate: nhpa}
	has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}

	/*
		instance := &resourcev1alpha1.HCPHybridAutoScaler{
			ObjectMeta: metav1.ObjectMeta{
				Name: cluster + "-" + hpa.Spec.ScaleTargetRef.Name + "-hpa2",
			},
			Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
				TargetCluster: cluster,
				WarningCount:  2,
				ScalingOptions: resourcev1alpha1.ScalingOptions{
					HpaTemplate: nhpa,
				},
			},
			Status: resourcev1alpha1.HCPHybridAutoScalerStatus{
				ResourceStatus: "WAITING",
			},
		}
	*/
	newhas, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("update %s Done\n", newhas.Name)
		return nil
	}
}

func (a *autoscaler) CreateVPA(deployment *appsv1.Deployment, updateMode string) error {

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

	// has_name := a.cluster + "-" + a.deployment.Name
	has, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), (*a.has_name)[deployment], metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 3. vpaTemplate -> HCPHybridAutoScaler 생성
	has.Status.LastSpec = has.Spec
	has.Spec.WarningCount = 3
	has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{VpaTemplate: vpa}
	has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}

	/*

		instance := &resourcev1alpha1.HCPHybridAutoScaler{
			ObjectMeta: metav1.ObjectMeta{
				Name: cluster + "-" + vpa.Name + "-vpa",
			},
			Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
				TargetCluster: cluster,
				WarningCount:  3,
				ScalingOptions: resourcev1alpha1.ScalingOptions{
					VpaTemplate: vpa,
				},
			},
			Status: resourcev1alpha1.HCPHybridAutoScalerStatus{
				ResourceStatus: "WAITING",
			},
		}
	*/
	newhas, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("update %s Done\n", newhas.Name)
		return nil
	}
}
