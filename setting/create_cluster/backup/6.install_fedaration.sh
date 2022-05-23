kubectl apply -f RBAC.yaml 
#helm init --service-account tiller
#helm init --service-account tiller --output yaml | sed 's@apiVersion: extensions/v1beta1@apiVersion: apps/v1@' | sed 's@  replicas: 1@  replicas: 1\n  selector: {"matchLabels": {"app": "helm", "name": "tiller"}}@' | kubectl apply -f -
kubectl apply -f tiller.yaml
helm repo add kubefed-charts https://raw.githubusercontent.com/kubernetes-sigs/kubefed/master/charts
helm repo list
helm search repo

kubectl create ns kube-federation-system
helm install kubefed kubefed-charts/kubefed --version=0.8.1 --namespace kube-federation-system
#helm install kubefed kubefed-charts/kubefed --version=0.1.0-rc6 --namespace kube-federation-system 
