package resource

import (
	"Hybrid_Cluster/hcp-scheduler/pkg/resource"
	cobrautil "Hybrid_Cluster/hybridctl/util"
	resourcev1alpha1 "Hybrid_Cluster/pkg/apis/resource/v1alpha1"
	hasv1alpha1 "Hybrid_Cluster/pkg/client/resource/v1alpha1/clientset/versioned"
	cm "Hybrid_Cluster/util/clusterManager"
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpav1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
)

func GetPod(cluster string, pod string, pod_namespace string) (*corev1.Pod, error) {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)
	p := &corev1.Pod{}
	p, err := cluster_client.CoreV1().Pods(pod_namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	if err != nil {
		return p, err
	} else {
		fmt.Printf("success to get pod %s [cluster %s, node %s]\n", p.Name, cluster, p.Spec.NodeName)
		return p, err
	}
}

func GetDeploymentName(pod *corev1.Pod) string {
	str := pod.GenerateName
	list := strings.SplitAfterN(str, "-", -1)
	length := len(list[len(list)-2]) + 1
	name := str[:len(str)-length]
	return name
}

func GetDeployment(cluster string, pod *corev1.Pod) (*v1.Deployment, error) {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)
	d := &v1.Deployment{}
	deploymentName := GetDeploymentName(pod)
	d, err := cluster_client.AppsV1().Deployments(pod.Namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return d, err
	} else {
		fmt.Printf("success to get deployment %s in cluster %s [replicas : %d]\n", d.Name, cluster, *d.Spec.Replicas)
	}
	return d, err
}

func CreateDeployment(cluster string, node string, deployment *v1.Deployment) error {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)
	deployment.Spec.Template.Spec.NodeName = node
	deployment.ResourceVersion = ""
	new_dep, err := cluster_client.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		return err
	} else {
		fmt.Printf("success to create %s in cluster %s [replicas : %d]\n", new_dep.Name, cluster, *deployment.Spec.Replicas)
	}
	return nil
}

func UpdateDeployment(cluster string, deployment *v1.Deployment, replicas int32) error {
	config, _ := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	cluster_client := kubernetes.NewForConfigOrDie(config)

	if deployment.Status.Replicas == deployment.Status.ReadyReplicas && deployment.Status.Replicas == deployment.Status.AvailableReplicas {
		fmt.Println("deployment status is OK")
		deployment.Spec.Replicas = &replicas
		new_dep, err := cluster_client.AppsV1().Deployments("default").Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		} else {
			fmt.Printf("success to update %s in cluster %s [replicas : %d]\n", new_dep.Name, cluster, *deployment.Spec.Replicas)
		}
	} else {
		fmt.Println("deployment status is not stable")
	}
	return nil
}

func CreateHPA(cluster string, pod string, namespace string, minReplicas *int32, maxReplicas int32) error {
	cm := cm.NewClusterManager()
	master_config := cm.Host_config
	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(master_config)
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
	hpa := hpav2beta1.HorizontalPodAutoscaler{
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
			Name: hpa.Spec.ScaleTargetRef.Name + "-hpa",
		},
		Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
			WarningCount: 1,
			CurrentStep:  "HAS", // HAS -> Sync -> Done
			ScalingOptions: resourcev1alpha1.ScalingOptions{
				HpaTemplate: hpa,
			},
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
	cm := cm.NewClusterManager()
	master_config := cm.Host_config

	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(master_config)
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
	client, err := kubernetes.NewForConfig(master_config)
	hpa, err := client.AutoscalingV2beta1().HorizontalPodAutoscalers(d.Namespace).Get(context.TODO(), d.Name, metav1.GetOptions{})
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

	instance := &resourcev1alpha1.HCPHybridAutoScaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: hpa.Spec.ScaleTargetRef.Name + "-hpa2",
		},
		Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
			WarningCount: 2,
			CurrentStep:  "HAS", // HAS -> Sync -> Done
			ScalingOptions: resourcev1alpha1.ScalingOptions{
				HpaTemplate: nhpa,
			},
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

func CreateVPA(cluster string, pod string, namespace string, updateMode string) error {
	cm := cm.NewClusterManager()
	master_config := cm.Host_config
	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(master_config)
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

	// 3. vpaTemplate -> HCPHybridAutoScaler 생성
	instance := &resourcev1alpha1.HCPHybridAutoScaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: vpa.Name + "-vpa",
		},
		Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
			WarningCount: 3,
			CurrentStep:  "HAS", // HAS -> Sync -> Done
			ScalingOptions: resourcev1alpha1.ScalingOptions{
				VpaTemplate: vpa,
			},
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
