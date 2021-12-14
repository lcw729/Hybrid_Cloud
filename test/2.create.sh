#!/bin/bash
kubectl create ns hcp
kubectl create -f deploy/operator.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/rolebinding.yaml


