package main

import (
	"flag"
	"fmt"
	"time"

	controller "github.com/KETI-Hybrid/hcp-has-controller-v1/src/controller"
	"github.com/KETI-Hybrid/hcp-pkg/util/clusterManager"

	v1alpha1hcphas "github.com/KETI-Hybrid/hcp-pkg/client/resource/v1alpha1/clientset/versioned"

	kubeinformers "k8s.io/client-go/informers"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"

	informers "github.com/KETI-Hybrid/hcp-pkg/client/resource/v1alpha1/informers/externalversions"

	"k8s.io/sample-controller/pkg/signals"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cm, err := clusterManager.NewClusterManager()
	if err != nil {
		fmt.Println(err)
		return
	}
	hcphasclient, err := v1alpha1hcphas.NewForConfig(cm.Host_config)
	if err != nil {
		klog.Info(err)
	}
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(cm.Host_kubeClient, time.Second*30)
	hcphasInformerFactory := informers.NewSharedInformerFactory(hcphasclient, time.Second*30)

	controller := controller.NewController(cm.Host_kubeClient, hcphasclient, hcphasInformerFactory.Hcp().V1alpha1().HCPHybridAutoScalers())
	kubeInformerFactory.Start(stopCh)
	hcphasInformerFactory.Start(stopCh)
	if err := controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
