package main

/*
Copyright 2018 The Multicluster-Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	controller "Hybrid_Cluster/hcp-has-controller/pkg/controller"
	"Hybrid_Cluster/util/clusterManager"
	"flag"
	"time"

	v1alpha1hcphas "Hybrid_Cluster/pkg/client/resource/v1alpha1/clientset/versioned"

	kubeinformers "k8s.io/client-go/informers"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"

	informers "Hybrid_Cluster/pkg/client/resource/v1alpha1/informers/externalversions"

	"k8s.io/sample-controller/pkg/signals"
)

// var c chan string

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cm := clusterManager.NewClusterManager()
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
