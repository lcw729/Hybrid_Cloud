#!/bin/bash
cd deploy


kubectl delete -f service_account.yaml
kubectl delete -f role_binding.yaml
kubectl delete -f operator.yaml
kubectl delete -f crds/cr.yaml

kubectl delete ns hybrid 
kubectl delete deploy example-hcppolicy-deploy --context cluster1 -n hybrid
kubectl delete deploy example-hcppolicy-deploy --context cluster2 -n hybrid
kubectl delete deploy example-hcppolicy-deploy --context cluster3 -n hybrid
kubectl delete hcppolicyengines example-hcppolicy -n hybrid
#kubectl delete ns hybrid --context cluster1 &
#kubectl delete ns hybrid --context cluster2 &
#kubectl delete ns hybrid --context cluster3 &
#kubectl delete ns hybrid &

kubectl delete -f crds/crd.yaml
cd ..
