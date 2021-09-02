#!/bin/bash
kubectl delete ns "kube-federation-system" --context gke_keti-container_us-central1-a_cluster-1
kubectl delete clusterroles.rbac.authorization.k8s.io kubefed-controller-manager:cluster-1 --context gke_keti-container_us-central1-a_cluster-1 
kubectl delete clusterrolebindings.rbac.authorization.k8s.io kubefed-controller-manager:cluster-1-hcp --context gke_keti-container_us-central1-a_cluster-1
kubectl delete kubefedclusters cluster-1 -n kube-federation-system --context kube-master