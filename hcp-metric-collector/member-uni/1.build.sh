#!/bin/bash

controller_name="hcp-metric-collector"
password="ketilinux"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor Hybrid_Cloud$controller_name/member-uni/cmd/main

# gsutil cp build/_output/bin/$controller_name gs://khg-bucket/

# kubectl create -f deploy/volume/ --context uni-master


sshpass -p $password scp -r build/_output/bin/$controller_name root@10.0.5.43:~/$controller_name-build


sshpass -p $password ssh root@10.0.5.43 "cd $controller_name-build; ./1.upload.sh"