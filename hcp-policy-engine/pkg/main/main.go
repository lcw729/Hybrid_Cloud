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

package main

import (
	// "Hybrid_Cloud/hcp-policy-engine/pkg/controller"

	"flag"
	"fmt"
	"time"

	"hcp-pkg/util/clusterManager"

	controller "hcp-policy-engine/pkg/controller"

	v1alpha1hcppolicy "hcp-pkg/client/hcppolicy/v1alpha1/clientset/versioned"

	kubeinformers "k8s.io/client-go/informers"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"

	informers "hcp-pkg/client/hcppolicy/v1alpha1/informers/externalversions"

	"k8s.io/sample-controller/pkg/signals"
)

// var c chan string

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cm, err := clusterManager.NewClusterManager()
	if err != nil {
		fmt.Println(err)
		return
	}
	hcppolicyclient, err := v1alpha1hcppolicy.NewForConfig(cm.Host_config)
	if err != nil {
		klog.Info(err)
	}
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(cm.Host_kubeClient, time.Second*30)
	hcppolicyInformerFactory := informers.NewSharedInformerFactory(hcppolicyclient, time.Second*30)

	controller := controller.NewController(cm.Host_kubeClient, hcppolicyclient, hcppolicyInformerFactory.Hcp().V1alpha1().HCPPolicies())
	kubeInformerFactory.Start(stopCh)
	hcppolicyInformerFactory.Start(stopCh)
	if err := controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
