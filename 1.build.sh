#!/bin/bash
docker_id="ketidevit2"
controller_name="hcp-metric-collector"

export GO111MODULE=on
go mod vendor

go build -o build/_output/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor Hybrid_Cloud/hcp-metric-collector/master/cmd/main && \

docker build -t $docker_id/$controller_name:v0.0.31 build && \
docker push $docker_id/$controller_name:v0.0.31


#0.0.23 >> 운용중인 metric collector
#0.0.31 >> test metric collector --> 제병 테스트