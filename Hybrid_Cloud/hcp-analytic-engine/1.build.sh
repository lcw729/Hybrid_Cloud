#!/bin/bash
docker_id="ketidevit2"
controller_name="hcp-analytic-engine"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor Hybrid_Cloud/hcp-github.com/hcp-analytic-engine-v1/pkg/main && \

docker build -t $docker_id/$controller_name:v0.0.1 build && \
docker push $docker_id/$controller_name:v0.0.1

