package main

import (
	controller "Hybrid_Cloud/hcp-scheduler/src/controller"
	"flag"
	"time"

	resourcev1alpha1 "Hybrid_Cloud/pkg/client/resource/v1alpha1/clientset/versioned"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	informers "Hybrid_Cloud/pkg/client/resource/v1alpha1/informers/externalversions"

	"k8s.io/sample-controller/pkg/signals"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	master_config, _ := rest.InClusterConfig()
	master_clientset, _ := kubernetes.NewForConfig(master_config)

	resourceclient, err := resourcev1alpha1.NewForConfig(master_config)
	if err != nil {
		klog.Info(err)
	}
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(master_clientset, time.Second*30)
	resourceInformerFactory := informers.NewSharedInformerFactory(resourceclient, time.Second*30)

	controller := controller.NewController(master_clientset, resourceclient, resourceInformerFactory.Hcp().V1alpha1().HCPDeployments())
	kubeInformerFactory.Start(stopCh)
	resourceInformerFactory.Start(stopCh)
	if err := controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
