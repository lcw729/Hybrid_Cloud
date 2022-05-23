# kubectl config use-context kube-master
# kubectl create ns hcp
kubectl apply -f crds/crd.yaml
kubectl apply -f crds/cr.yaml
kubectl apply -f policy/optimalArragementAlgorithm.yaml
kubectl apply -f policy/resourceConfigurationCycle.yaml
kubectl apply -f policy/warningLevel.yaml
kubectl apply -f policy/watchingLevel.yaml
kubectl apply -f policy/weightCalculationCycle.yaml
