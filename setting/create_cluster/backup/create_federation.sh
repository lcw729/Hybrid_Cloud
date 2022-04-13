# Install My Federation
echo "---------------------"
echo "Install My Federation"
echo "---------------------"
helm install /root/workspace/go/src/sigs.k8s.io/kubefed/charts/kubefed --name kubefed --version=0.1.0-rc6 --namespace kube-federation-system

