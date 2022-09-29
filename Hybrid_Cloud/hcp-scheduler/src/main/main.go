package main

import (
	"flag"
	"time"

	"github.com/KETI-Hybrid/hcp-pkg/util/clusterManager"

	controller "github.com/KETI-Hybrid/hcp-scheduler-v1/src/controller"

	kubeinformers "k8s.io/client-go/informers"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"

	informers "github.com/KETI-Hybrid/hcp-pkg/client/resource/v1alpha1/informers/externalversions"

	"k8s.io/sample-controller/pkg/signals"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	cm, err := clusterManager.NewClusterManager()
	if err != nil {
		klog.Errorln(err)
	}

	stopCh := signals.SetupSignalHandler()
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(cm.Host_kubeClient, time.Second*30)
	resourceInformerFactory := informers.NewSharedInformerFactory(cm.HCPResource_Client, time.Second*30)

	controller := controller.NewController(cm.Host_kubeClient, cm.HCPResource_Client, resourceInformerFactory.Hcp().V1alpha1().HCPDeployments())
	kubeInformerFactory.Start(stopCh)
	resourceInformerFactory.Start(stopCh)
	if err := controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
