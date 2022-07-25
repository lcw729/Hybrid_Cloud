package handler

import (
	resourcev1alpha1 "Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

type Resource struct {
	TargetCluster string
	RealResource  interface{}
}

func CreateDeploymentHandler(w http.ResponseWriter, r *http.Request) {

	var resource Resource
	jsonDataFromHttp, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(jsonDataFromHttp, &resource)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	// RealResource 읽어오기
	var real_resource *appsv1.Deployment
	bytes, _ := json.Marshal(resource.RealResource)
	json.Unmarshal(bytes, &real_resource)

	cm, _ := clusterManager.NewClusterManager()
	// TargetCluster가 지정되지 않은 경우
	if resource.TargetCluster == "" {

		// HCPDeployment 생성하기
		hcp_resource := deploymentToHCPDeployment(real_resource)

		r, err := cm.HCPResource_Client.HcpV1alpha1().HCPDeployments("hcp").Create(context.TODO(), &hcp_resource, metav1.CreateOptions{})
		if err != nil {
			klog.Error(err)
			return
		} else {
			klog.Info("Request scheduling to scheduler : %s \n", r.Name)
		}
	} else {
		// TargetCluster가 지정된 경우
		target_clientset := cm.Cluster_kubeClients[resource.TargetCluster]
		// namespace
		namespace := real_resource.ObjectMeta.Namespace
		if namespace == "" {
			namespace = "default"
		}

		hcp_resource := deploymentToHCPDeployment(real_resource)
		hcp_resource.Spec.SchedulingResult.Targets = append(hcp_resource.Spec.SchedulingResult.Targets, resourcev1alpha1.Target{
			Cluster:  resource.TargetCluster,
			Replicas: real_resource.Spec.Replicas,
		})
		_, err = cm.HCPResource_Client.HcpV1alpha1().HCPDeployments("hcp").Create(context.TODO(), &hcp_resource, metav1.CreateOptions{})
		if err != nil {
			klog.Error(err)
			return
		} else {
			klog.Info("Succeed to create hcpdeployment: %s \n", hcp_resource.Name)
		}
		// Kubernetes Deployment 생성
		r, err := target_clientset.AppsV1().Deployments(namespace).Create(context.TODO(), real_resource, metav1.CreateOptions{})

		if err != nil {
			klog.Error(err)
			return
		} else {
			klog.Info("Succeed to create deployment %s \n", r.Name)
		}
	}
}

func DeleteDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	hcpdeployment, err := cm.HCPResource_Client.HcpV1alpha1().HCPDeployments("hcp").Get(context.TODO(), name, metav1.GetOptions{})

	if !hcpdeployment.Spec.SchedulingNeed && hcpdeployment.Spec.SchedulingComplete {
		// if target_cluster != "" {
		// 	config, _ := cobrautil.BuildConfigFromFlags(target_cluster, "/root/.kube/config")
		// 	clientset, _ := kubernetes.NewForConfig(config)
		// 	err = clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		// } else {
		targets := hcpdeployment.Spec.SchedulingResult.Targets
		for _, target := range targets {
			// TODO : cluster unregister한 경우
			cm, _ := clusterManager.NewClusterManager()
			clientset := cm.Cluster_kubeClients[target.Cluster]
			err = clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		}
		// }
	}

	if err != nil {
		klog.Error(err)
		return
	} else {
		klog.Info("Succeed to delete deployment %s \n", name)
		err = cm.HCPResource_Client.HcpV1alpha1().HCPDeployments("hcp").Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			klog.Error(err)
			return
		} else {
			klog.Info("Succeed to delete hcpdeployment %s \n", name)
		}
	}
}

func CreatePodHandler(w http.ResponseWriter, r *http.Request) {

	var resource Resource

	jsonDataFromHttp, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(jsonDataFromHttp, &resource)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	// RealResource 읽어오기
	var real_resource *v1.Pod
	bytes, _ := json.Marshal(resource.RealResource)
	json.Unmarshal(bytes, &real_resource)

	// TargetCluster가 지정되지 않은 경우
	if resource.TargetCluster == "undefined" {
		// HCPDeployment 생성하기
		hcp_resource := resourcev1alpha1.HCPPod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "HCPPod",
				APIVersion: "hcp.crd.com/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: real_resource.Name,
			},
			Spec: resourcev1alpha1.HCPPodSpec{
				RealPodSpec:     real_resource.Spec,
				RealPodMetadata: real_resource.ObjectMeta,

				// SchedulingStatus "Requested"
				SchedulingStatus: "Requested",
				// SchedulingType:   algorithm,
			},
		}

		r, err := cm.HCPResource_Client.HcpV1alpha1().HCPPods("hcp").Create(context.TODO(), &hcp_resource, metav1.CreateOptions{})
		if err != nil {
			klog.Error(err)
			return
		} else {
			klog.Info("request scheduling to scheduler : %s \n", r.Name)
		}
	} else {
		// TargetCluster가 지정된 경우
		cm, _ := clusterManager.NewClusterManager()
		clientset := cm.Cluster_kubeClients[resource.TargetCluster]
		// namespace
		namespace := real_resource.ObjectMeta.Namespace
		if namespace == "" {
			namespace = "default"
		}

		// Kubernetes Deployment 생성
		r, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), real_resource, metav1.CreateOptions{})

		if err != nil {
			klog.Error(err)
			return
		} else {
			klog.Info("Succeed to create pod %s \n", r.Name)
		}
	}
}

func CreateHCPHASHandler() {

}

func deploymentToHCPDeployment(real_resource *appsv1.Deployment) resourcev1alpha1.HCPDeployment {
	// HCPDeployment 생성하기
	hcp_resource := resourcev1alpha1.HCPDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HCPDeployment",
			APIVersion: "hcp.crd.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: real_resource.Name,
		},
		Spec: resourcev1alpha1.HCPDeploymentSpec{
			RealDeploymentSpec:     real_resource.Spec,
			RealDeploymentMetadata: real_resource.ObjectMeta,

			// SchedulingStatus "Requested"
			SchedulingNeed:     true,
			SchedulingComplete: false,
			//SchedulingType:   algorithm[0],
		},
	}
	return hcp_resource
}
