#!/bin/bash

# install kubefed
./create_cluster/4.get_kubefedctl.sh
./create_cluster/5.install_helm.sh
./create_cluster/6.install_fedaration.sh

# create CRD
kubectl apply -f /root/go/src/KETI-Hybrid_Cloud/pkg/crds

