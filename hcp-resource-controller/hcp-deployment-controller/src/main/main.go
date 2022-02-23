package main

import (
	controller "Hybrid_Cluster/hcp-resource-controller/hcp-deployment-controller/src/controller"
	"Hybrid_Cluster/util/clusterManager"
	"flag"
	"time"

	resourcev1alpha1 "Hybrid_Cluster/pkg/client/resource/v1alpha1/clientset/versioned"

	kubeinformers "k8s.io/client-go/informers"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"

	informers "Hybrid_Cluster/pkg/client/resource/v1alpha1/informers/externalversions"

	"k8s.io/sample-controller/pkg/signals"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cm := clusterManager.NewClusterManager()
	client, err := resourcev1alpha1.NewForConfig(cm.Host_config)
	if err != nil {
		klog.Info(err)
	}
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(cm.Host_kubeClient, time.Second*30)
	resourcev1alpha1InformerFactory := informers.NewSharedInformerFactory(client, time.Second*30)

	controller := controller.NewController(cm.Host_kubeClient, client, resourcev1alpha1InformerFactory.Hcp().V1alpha1().HCPDeployments())
	kubeInformerFactory.Start(stopCh)
	resourcev1alpha1InformerFactory.Start(stopCh)
	if err := controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
