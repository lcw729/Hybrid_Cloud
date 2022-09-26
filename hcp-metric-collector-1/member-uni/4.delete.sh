#!/bin/bash


# kubectl delete -f deploy/ --context uni-master

kubectl delete -f deploy/operator/operator-uni-cluster.yaml --context uni-master


# kubectl delete -f deploy/secret.yaml --context uni-master
kubectl delete -f deploy/service.yaml --context uni-master
kubectl delete -f deploy/service_account_before.yaml --context uni-master
kubectl delete -f deploy/role_binding.yaml --context uni-master



kubectl delete -f deploy/volume/pvc_fedora.yml --context uni-master
# kubectl delete -f deploy/volume/pvc_sample.yaml --context uni-master
# kubectl delete -f deploy/volume/pv_sample.yaml --context uni-master
kubectl delete -f deploy/volume/pv_fedora.yml --context uni-master

