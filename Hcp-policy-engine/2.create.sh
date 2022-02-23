#!/bin/bash
cd deploy

kubectl create ns hybrid
kubectl create ns hybrid
# kubectl create ns hybrid --context cluster1
# kubectl create ns hybrid --context cluster2
# kubectl create ns hybrid --context cluster3
kubectl create -f crds/crd.yaml
kubectl create -f service_account.yaml
kubectl create -f role_binding.yaml
kubectl create -f operator.yaml
kubectl create -f crds/cr.yaml

cd ..
