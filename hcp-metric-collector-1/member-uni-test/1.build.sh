#!/bin/bash

controller_name="hcp-metric-collector"
password="ketilinux"
uni_server="10.0.5.43"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor Hybrid_Cloud/$controller_name/member-uni-test/cmd/main
# go build -o build/_output/bin/hcp-metric-collector -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor Hybrid_Cloud/hcp-metric-collector/member-uni-test/cmd/main

# gsutil cp build/_output/bin/$controller_name gs://khg-bucket/

# kubectl create -f deploy/volume/ --context uni-master


sshpass -p $password scp -r build/_output/bin/$controller_name root@$uni_server:~/$controller_name-build


sshpass -p $password ssh root@$uni_server "cd $controller_name-build; ./1.upload.sh"
