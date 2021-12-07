VERSION=0.8.1
OS=linux
ARCH=amd64
curl -LO https://github.com/kubernetes-sigs/kubefed/releases/download/v${VERSION}/kubefedctl-${VERSION}-${OS}-${ARCH}.tgz
tar -xvzf kubefedctl-0.8.1-linux-amd64.tgz
chmod u+x kubefedctl
sudo mv kubefedctl /usr/local/bin/ #make sure the location is in the PATH

# kubefedctl version

# curl -LO https://git.io/get_helm.sh
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh
helm repo update


# helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kubefed-charts https://raw.githubusercontent.com/kubernetes-sigs/kubefed/master/charts
helm repo list
helm search repo
kubectl create ns kube-federation-system
helm --namespace kube-federation-system upgrade -i kubefed kubefed-charts/kubefed --version=0.8.1
kubectl get pod  -n kube-federation-system