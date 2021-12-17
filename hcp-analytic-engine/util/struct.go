package util

type Cluster struct {
	ClusterInfo *ClusterInfo
	Nodes       []*NodeScore
}

type NodeScore struct {
	NodeId int32 `protobuf:"varint,1,opt,name=node_id" json:"node_id,omitempty"`
	Score  int32 `protobuf:"varint,2,opt,name=score" json:"score,omitempty"`
}

type ClusterInfo struct {
	ClusterId   int32  `protobuf:"varint,1,opt,name=cluster_id" json:"cluster_id,omitempty"`
	ClusterName string `protobuf:"bytes,2,opt,name=cluster_name" json:"cluster_name,omitempty"`
}

type WatchingLevel struct {
	Levels []Level
}

type Level struct {
	Type  string
	Value string
}
