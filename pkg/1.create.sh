# kubectl apply -f crds/clusterregister.yaml
kubectl apply -f crds/vpa/vpa-v1-crd-gen.yaml
kubectl apply -f crds/hcpcluster.yaml
kubectl apply -f crds/sync.yaml
kubectl apply -f crds/has.yaml
# kubectl apply -f crds/kubefedcluster.yaml