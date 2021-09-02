#!/bin/bash
cd deploy


kubectl delete -f service_account.yaml
kubectl delete -f role_binding.yaml
kubectl delete -f operator.yaml
kubectl delete -f crds/cr.yaml

kubectl delete deploy example-templateresource-deploy --context cluster1 -n nsnsns
kubectl delete deploy example-templateresource-deploy --context cluster2 -n nsnsns
kubectl delete deploy example-templateresource-deploy --context cluster3 -n nsnsns
kubectl delete templateresources example-templateresource -n nsnsns
#kubectl delete ns nsnsns --context cluster1 &
#kubectl delete ns nsnsns --context cluster2 &
#kubectl delete ns nsnsns --context cluster3 &
#kubectl delete ns nsnsns &

kubectl delete -f crds/crd.yaml
cd ..
