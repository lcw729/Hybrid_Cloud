#!/bin/bash


kubectl delete -f deploy/ --context uni-master


kubectl delete -f deploy/operator/operator-test-cluster.yaml --context uni-master

# kubectl delete -f deploy/volum --context uni-master