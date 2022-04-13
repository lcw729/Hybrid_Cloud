helm repo add kubefed-charts https://raw.githubusercontent.com/kubernetes-sigs/kubefed/master/charts
helm repo list
helm search repo

kubectl create ns kube-federation-system
helm install kubefed kubefed-charts/kubefed --version=0.8.1 --namespace kube-federation-system
