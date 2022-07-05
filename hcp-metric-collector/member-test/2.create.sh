#!/bin/bash


kubectl create ns hcp --context uni-master
kubectl create -f deploy/ --context uni-master
kubectl create -f deploy/operator/operator-test-cluster.yaml --context uni-master
