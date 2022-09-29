# Delete all Running Federation Object
echo "------------------------------------"
echo "Delete All Running Federation Object"
echo "------------------------------------"
#kubectl delete daemonsets,replicasets,services,deployments,pods,rc -n kube-federation-system --all
kubectl delete all -n kube-federation-system --all

# Unjoin Clusters
echo "-------------------------------------"
echo "           Unjoun Clusters           "
echo "-------------------------------------"
kubefedctl unjoin cluster2 --cluster-context cluster2 --host-cluster-context openmcp --v=2
kubefedctl unjoin cluster1 --cluster-context cluster1 --host-cluster-context openmcp --v=2

# Delete FederatedTypeConfig
echo "-------------------------------------"
echo "     Delete FederatedTypeConfig      "
echo "-------------------------------------"
kubectl -n kube-federation-system delete FederatedTypeConfig --all

# Delete CRD
echo "-------------------------------------"
echo "             Delete CRD              " 
echo "-------------------------------------"
kubectl delete crd $(kubectl get crd | grep -E 'kubefed.io' | awk '{print $1}')

# Delete Federation for helm
echo "-------------------------------------"
echo "     Delete Federation for helm      "
echo "-------------------------------------"
helm delete --purge kubefed
