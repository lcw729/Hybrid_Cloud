# Join Clusters
echo "---------------------"
echo "    Join Clusters    "
echo "---------------------"
kubefedctl join cluster1 --cluster-context cluster1 --host-cluster-context openmcp --v=2
kubefedctl join cluster2 --cluster-context cluster2 --host-cluster-context openmcp --v=2

# Check Cluster Combine
echo "---------------------"
echo "Check Cluster Combine"
echo "---------------------"
kubectl -n kube-federation-system get kubefedclusters
