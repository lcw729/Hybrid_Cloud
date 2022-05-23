swapoff -a
kubeadm init --pod-network-cidr=10.244.0.0/16 --ignore-preflight-errors=NumCPU --v=4 
#--cri-socket /var/run/crio/crio.sock 

mkdir -p $HOME/.kube
cp /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

sysctl net.bridge.bridge-nf-call-iptables=1
#kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/bc79dd1505b0c8681ece4de4c0d86c5cd2643275/Documentation/kube-flannel.yml
#kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/62e44c867a2846fefb68bd5f178daf4da3095ccb/Documentation/kube-flannel.yml
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
#kubectl apply -f calico.yaml
