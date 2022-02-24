package resource

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"

	cobrautil "Hybrid_Cloud/hybridctl/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// replicas := int32(3)

	// deployment := v1.Deployment{
	// 	TypeMeta: metav1.TypeMeta{
	// 		Kind:       "Deployment",
	// 		APIVersion: "apps/v1",
	// 	},
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      "nginx-deployment",
	// 		Namespace: "default",
	// 		Labels:    map[string]string{"app": "nginx"},
	// 	},
	// 	Spec: v1.DeploymentSpec{
	// 		Replicas: &replicas,
	// 		Selector: &metav1.LabelSelector{
	// 			MatchLabels: map[string]string{
	// 				"app": "nginx",
	// 			},
	// 		},
	// 		Template: corev1.PodTemplateSpec{
	// 			ObjectMeta: metav1.ObjectMeta{
	// 				Labels: map[string]string{"app": "nginx"},
	// 			},
	// 			Spec: corev1.PodSpec{
	// 				Containers: []corev1.Container{
	// 					{
	// 						Name:  "nginx",
	// 						Image: "nginx:1.14.2",
	// 						Ports: []corev1.ContainerPort{
	// 							{
	// 								ContainerPort: 80,
	// 							},
	// 						},
	// 					},
	// 				},
	// 				NodeName: "cluster2-worker2",
	// 			},
	// 		},
	// 	},
	// }
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
