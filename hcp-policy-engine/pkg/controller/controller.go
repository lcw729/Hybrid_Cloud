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

package controller

import (
	"Hybrid_Cluster/apis"
	policyv1alpha1 "Hybrid_Cluster/apis/policy/v1alpha1"
	"Hybrid_Cluster/hcplog"
	"Hybrid_Cluster/util/clusterManager"
	"context"
	"fmt"

	resourcev1alpha1 "Hybrid_Cluster/apis/resource/v1alpha1"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller_hcphybridautoscaler")

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	hcplog.V(4).Info("[HCP Policy Engine] Function Called NewController")
	cm = myClusterManager
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	hcplog.V(4).Info(live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(context.Background(), live, &policyv1alpha1.HCPPolicy{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	i += 1
	hcplog.V(4).Info("********* [", i, "] *********")
	hcplog.V(5).Info("Request Context: ", req.Context, " / Request Namespace: ", req.Namespace, " /  Request Name: ", req.Name)

	// Fetch the HCPDeployment instance
	instance := &policyv1alpha1.HCPPolicy{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	hcplog.V(5).Info("instance Name: ", instance.Name)
	hcplog.V(5).Info("instance Namespace: ", instance.Namespace)

	if err != nil {
		if errors.IsNotFound(err) {
			hcplog.V(2).Info("Delete Policy Resource")
			return reconcile.Result{}, nil
		}
		hcplog.V(0).Info("Error: ", err)

		return reconcile.Result{}, err
	}

	if instance.Spec.PolicyStatus == "Disabled" {
		hcplog.V(2).Info("Policy Disabled")
	} else if instance.Spec.PolicyStatus == "Enabled" {
		if instance.Spec.RangeOfApplication == "FromNow" {
			hcplog.V(2).Info("Policy Enabled - FromNow")
		} else if instance.Spec.RangeOfApplication == "All" {
			object := instance.Spec.Template.Spec.TargetController.Kind
			if object == "HCPHybridAutoScaler" {
				hcplog.V(2).Info("Policy Enabled - HCPHybridAutoScaler")
				hpaList := &resourcev1alpha1.HCPHybridAutoScalerList{}
				listOptions := &client.ListOptions{Namespace: ""} //all resources
				r.live.List(context.TODO(), hpaList, listOptions)
				for _, hpaInstance := range hpaList.Items {
					var i = 0
					for index, tmpPolicy := range hpaInstance.Status.Policies { //Find target policy
						if tmpPolicy.Type == instance.Spec.Template.Spec.Policies[0].Type { //Already exists
							i++
							hpaInstance.Status.Policies[index].Value = instance.Spec.Template.Spec.Policies[0].Value
							break
						}
					}
					if i == 0 {
						hpaInstance.Status.Policies = append(hpaInstance.Status.Policies, instance.Spec.Template.Spec.Policies...)
					}
					err := r.live.Status().Update(context.TODO(), &hpaInstance)
					if err != nil {
						hcplog.V(0).Info("HCPHPA Policy Update Error")
						return reconcile.Result{}, err
					} else {
						hcplog.V(2).Info("HCPHPA Policy UPDATE Success!")
					}
				}
			} else if object == "HCPLoadbalancer" {

			}
		}
	}
	return reconcile.Result{}, nil
}
