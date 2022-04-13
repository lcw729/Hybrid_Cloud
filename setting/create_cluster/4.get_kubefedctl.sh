VERSION=0.8.1
OS=linux
ARCH=amd64
curl -LO https://github.com/kubernetes-sigs/kubefed/releases/download/v${VERSION}/kubefedctl-${VERSION}-${OS}-${ARCH}.tgz
tar -zxvf kubefedctl-*.tgz
chmod u+x kubefedctl
mv kubefedctl /usr/local/bin/
