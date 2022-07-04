package main

import (
	controller "Hybrid_Cloud/hcp-cluster-manager/src/controller"
	hcpclusterv1alpha1 "Hybrid_Cloud/pkg/client/hcpcluster/v1alpha1/clientset/versioned"
	informers "Hybrid_Cloud/pkg/client/hcpcluster/v1alpha1/informers/externalversions"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"flag"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/sample-controller/pkg/signals"
)

func main() {
	config, err := rest.InClusterConfig()
	// config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Print config")
	fmt.Println(config)
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Print clientset")
	fmt.Println(clientset)
	// podsInNode, _ := hostKubeClient.CoreV1().Pods("").List(metav1.ListOptions{})
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	for _, p := range pods.Items {
		fmt.Println(p.GetName())
	}
	if err != nil {
		panic(err.Error())
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	for _, n := range nodes.Items {
		fmt.Println(n.GetName())
	}
	if err != nil {
		panic(err.Error())
	}
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cm, err := clusterManager.NewClusterManager()
	if err != nil {
		fmt.Println(err)
		return
	}
	hcpcluster_client, err := hcpclusterv1alpha1.NewForConfig(cm.Host_config)
	if err != nil {
		klog.Info(err)
	}
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(cm.Host_kubeClient, time.Second*30)
	hcpclusterInformerFactory := informers.NewSharedInformerFactory(hcpcluster_client, time.Second*30)
	//
	controller := controller.NewController(cm.Host_kubeClient, hcpcluster_client, hcpclusterInformerFactory.Hcp().V1alpha1().HCPClusters())
	kubeInformerFactory.Start(stopCh)
	hcpclusterInformerFactory.Start(stopCh)
	if err := controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}

}
