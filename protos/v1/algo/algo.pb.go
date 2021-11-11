// Code generated by protoc-gen-go.
// source: protos/v1/algo/algo.proto
// DO NOT EDIT!

/*
Package Hybrid_Cluster_protos_v1_algo is a generated protocol buffer package.

It is generated from these files:
	protos/v1/algo/algo.proto

It has these top-level messages:
	WeightResult
	ClusterInfo
	NodeScore
	Cluster
	ClusterWeightCalculatorRequest
	ClusterWeightCalculatorResponse
	OptimalArrangementRequest
	OptimalArrangementResponse
*/
package algo

import proto "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type WeightResult struct {
	ClusterId     int32  `protobuf:"varint,1,opt,name=cluster_id" json:"cluster_id,omitempty"`
	ClusterName   string `protobuf:"bytes,2,opt,name=cluster_name" json:"cluster_name,omitempty"`
	ClusterWeight int32  `protobuf:"varint,3,opt,name=cluster_weight" json:"cluster_weight,omitempty"`
}

func (m *WeightResult) Reset()         { *m = WeightResult{} }
func (m *WeightResult) String() string { return proto.CompactTextString(m) }
func (*WeightResult) ProtoMessage()    {}

type ClusterInfo struct {
	ClusterId   int32  `protobuf:"varint,1,opt,name=cluster_id" json:"cluster_id,omitempty"`
	ClusterName string `protobuf:"bytes,2,opt,name=cluster_name" json:"cluster_name,omitempty"`
}

func (m *ClusterInfo) Reset()         { *m = ClusterInfo{} }
func (m *ClusterInfo) String() string { return proto.CompactTextString(m) }
func (*ClusterInfo) ProtoMessage()    {}

type NodeScore struct {
	NodeId int32 `protobuf:"varint,1,opt,name=node_id" json:"node_id,omitempty"`
	Score  int32 `protobuf:"varint,2,opt,name=score" json:"score,omitempty"`
}

func (m *NodeScore) Reset()         { *m = NodeScore{} }
func (m *NodeScore) String() string { return proto.CompactTextString(m) }
func (*NodeScore) ProtoMessage()    {}

type Cluster struct {
	ClusterInfo *ClusterInfo `protobuf:"bytes,1,opt,name=cluster_info" json:"cluster_info,omitempty"`
	NodeScore   []*NodeScore `protobuf:"bytes,2,rep,name=node_score" json:"node_score,omitempty"`
}

func (m *Cluster) Reset()         { *m = Cluster{} }
func (m *Cluster) String() string { return proto.CompactTextString(m) }
func (*Cluster) ProtoMessage()    {}

func (m *Cluster) GetClusterInfo() *ClusterInfo {
	if m != nil {
		return m.ClusterInfo
	}
	return nil
}

func (m *Cluster) GetNodeScore() []*NodeScore {
	if m != nil {
		return m.NodeScore
	}
	return nil
}

// client - scheduler
type ClusterWeightCalculatorRequest struct {
}

func (m *ClusterWeightCalculatorRequest) Reset()         { *m = ClusterWeightCalculatorRequest{} }
func (m *ClusterWeightCalculatorRequest) String() string { return proto.CompactTextString(m) }
func (*ClusterWeightCalculatorRequest) ProtoMessage()    {}

// server - AnalyticEngine
type ClusterWeightCalculatorResponse struct {
	WeightResult []*WeightResult `protobuf:"bytes,1,rep,name=weight_result" json:"weight_result,omitempty"`
}

func (m *ClusterWeightCalculatorResponse) Reset()         { *m = ClusterWeightCalculatorResponse{} }
func (m *ClusterWeightCalculatorResponse) String() string { return proto.CompactTextString(m) }
func (*ClusterWeightCalculatorResponse) ProtoMessage()    {}

func (m *ClusterWeightCalculatorResponse) GetWeightResult() []*WeightResult {
	if m != nil {
		return m.WeightResult
	}
	return nil
}

// client - scheduler
type OptimalArrangementRequest struct {
}

func (m *OptimalArrangementRequest) Reset()         { *m = OptimalArrangementRequest{} }
func (m *OptimalArrangementRequest) String() string { return proto.CompactTextString(m) }
func (*OptimalArrangementRequest) ProtoMessage()    {}

// server - AnalyticEngine
type OptimalArrangementResponse struct {
	Status  bool       `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Cluster *Cluster   `protobuf:"bytes,2,opt,name=cluster" json:"cluster,omitempty"`
	Node    *NodeScore `protobuf:"bytes,3,opt,name=node" json:"node,omitempty"`
}

func (m *OptimalArrangementResponse) Reset()         { *m = OptimalArrangementResponse{} }
func (m *OptimalArrangementResponse) String() string { return proto.CompactTextString(m) }
func (*OptimalArrangementResponse) ProtoMessage()    {}

func (m *OptimalArrangementResponse) GetCluster() *Cluster {
	if m != nil {
		return m.Cluster
	}
	return nil
}

func (m *OptimalArrangementResponse) GetNode() *NodeScore {
	if m != nil {
		return m.Node
	}
	return nil
}

func init() {
}
