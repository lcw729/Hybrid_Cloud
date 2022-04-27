package resource

import (
	"Hybrid_Cloud/hcp-scheduler/backup/resource"
	resourcev1alpha1 "Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	hasv1alpha1 "Hybrid_Cloud/pkg/client/resource/v1alpha1/clientset/versioned"
	cm "Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"

	autoscaling "k8s.io/api/autoscaling/v1"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpav1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	"k8s.io/client-go/kubernetes"
)

func CreateHPA(cluster string, pod string, namespace string, minReplicas *int32, maxReplicas int32) error {
	cm, err := cm.NewClusterManager()
	if err != nil {
		fmt.Println(err)
		return err
	}
	config := cm.Host_config
	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 1. Pod 정보 -> Deployment 정보 얻기
	p, err := resource.GetPod(cluster, pod, namespace)
	if err != nil {
		fmt.Println(err)
		return err
	}
	d, err := resource.GetDeployment(cluster, p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 2. hapTemplate 생성
	hpa := &hpav2beta1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
		},
		Spec: hpav2beta1.HorizontalPodAutoscalerSpec{
			MinReplicas: minReplicas,
			MaxReplicas: maxReplicas,
			ScaleTargetRef: hpav2beta1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       d.Name,
			},
		},
	}

	// 3. hpaTemplate -> HCPHybridAutoScaler 생성
	instance := &resourcev1alpha1.HCPHybridAutoScaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: cluster + "-" + hpa.Spec.ScaleTargetRef.Name,
		},
		Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
			TargetCluster: cluster,
			WarningCount:  1,
			ScalingOptions: resourcev1alpha1.ScalingOptions{
				HpaTemplate: *hpa,
			},
		},
		Status: resourcev1alpha1.HCPHybridAutoScalerStatus{
			ResourceStatus: "WAITING",
		},
	}

	newhas, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Create(context.TODO(), instance, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("create %s Done\n", newhas.Name)
		return nil
	}
}

func CreateHPA2(cluster string, pod string, namespace string, minReplicas *int32, maxReplicas int32) error {
	cm, err := cm.NewClusterManager()
	if err != nil {
		return err
	}
	config := cm.Host_config
	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 1. Pod 정보 -> Deployment 정보 얻기
	p, err := resource.GetPod(cluster, pod, namespace)
	if err != nil {
		fmt.Println(err)
		return err
	}
	d, err := resource.GetDeployment(cluster, p)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 2. hapTemplate (warningCount 1) 정보 얻기
	target_cluster := cm.Cluster_configs[cluster]
	target_clientset, err := kubernetes.NewForConfig(target_cluster)
	if err != nil {
		fmt.Println(err)
		return err
	}

	hpa, err := target_clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(d.Namespace).Get(context.TODO(), d.Name, metav1.GetOptions{})
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
			MaxReplicas:    maxReplicas,
			ScaleTargetRef: hpa.Spec.ScaleTargetRef,
		},
	}

	has_name := cluster + "-" + d.Name
	has, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), has_name, metav1.GetOptions{})
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
	newhas, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("update %s Done\n", newhas.Name)
		return nil
	}
}

func CreateVPA(cluster string, pod string, namespace string, updateMode string) error {
	cm, err := cm.NewClusterManager()
	if err != nil {
		return err
	}
	config := cm.Host_config
	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 1. Pod 정보 -> Deployment 정보 얻기
	p, err := resource.GetPod(cluster, pod, namespace)
	if err != nil {
		fmt.Println(err)
		return err
	}
	d, err := resource.GetDeployment(cluster, p)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 2. vpaTemplate 생성
	// updateMode := vpav1beta2.UpdateModeAuto
	vpa := vpav1beta2.VerticalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "autoscaling.k8s.io/v1",
			Kind:       "VerticalPodAutoscaler",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
		},
		Spec: vpav1beta2.VerticalPodAutoscalerSpec{
			TargetRef: &autoscaling.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       d.Name,
			},
			UpdatePolicy: &vpav1beta2.PodUpdatePolicy{
				UpdateMode: (*vpav1beta2.UpdateMode)(&updateMode),
			},
		},
	}

	has_name := cluster + "-" + d.Name
	has, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), has_name, metav1.GetOptions{})
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
	newhas, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("update %s Done\n", newhas.Name)
		return nil
	}
}
