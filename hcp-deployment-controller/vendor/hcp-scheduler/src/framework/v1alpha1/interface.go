package v1alpha1

import (
	"hcp-scheduler/src/resourceinfo"
	"hcp-scheduler/src/util"

	v1 "k8s.io/api/core/v1"
)

type PluginScoreList map[string]util.TmpEachScore // key : plugin, value : 각 플러그인에 대한 각 클러스터의 Score 점수

type HCPFramework interface {
	//PluginScoreList{}
	RunFilterPluginsOnClusters(algorithms []string, pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfoList)
	RunScorePluginsOnClusters(algorithms []string, pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfoList)
}

type HCPPlugin interface {
	Name() string
}

type HCPFilterPlugin interface {
	HCPPlugin
	Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool
}

type HCPPostFilterPlugin interface {
	HCPPlugin
	PostFilter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo)
}

type HCPScorePlugin interface {
	HCPPlugin
	Score(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) int64
	Normalize(tmpEachScore *util.TmpEachScore, clusterInfoList *resourceinfo.ClusterInfoList)
}

/*
// Code

type Code int

const (
	Success Code = iota
	Error
	Unschedulable
	Wait
	Skip
)


// This list should be exactly the same as the codes iota defined above in the same order.
var codes = []string{"Success", "Error", "Unschedulable", "UnschedulableAndUnresolvable", "Wait", "Skip"}

// statusPrecedence defines a map from status to its precedence, larger value means higher precedent.
var statusPrecedence = map[Code]int{
	Error:         3,
	Unschedulable: 1,
	// Any other statuses we know today, `Skip` or `Wait`, will take precedence over `Success`.
	Success: -1,
}

// Status

type Status struct {
	code         Code
	reasons      []string
	err          error
	failedPlugin string
}

// Code returns code of the Status.
func (s *Status) Code() Code {
	if s == nil {
		return Success
	}
	return s.code
}

// Message returns a concatenated message on reasons of the Status.
func (s *Status) Message() string {
	if s == nil {
		return ""
	}
	return strings.Join(s.reasons, ", ")
}

// SetFailedPlugin sets the given plugin name to s.failedPlugin.
func (s *Status) SetFailedPlugin(plugin string) {
	s.failedPlugin = plugin
}

// FailedPlugin returns the failed plugin name.
func (s *Status) FailedPlugin() string {
	return s.failedPlugin
}

// Reasons returns reasons of the Status.
func (s *Status) Reasons() []string {
	return s.reasons
}

// AppendReason appends given reason to the Status.
func (s *Status) AppendReason(reason string) {
	s.reasons = append(s.reasons, reason)
}

// IsSuccess returns true if and only if "Status" is nil or Code is "Success".
func (s *Status) IsSuccess() bool {
	return s.Code() == Success
}

// IsUnschedulable returns true if "Status" is Unschedulable (Unschedulable or UnschedulableAndUnresolvable).
func (s *Status) IsUnschedulable() bool {
	code := s.Code()
	return code == Unschedulable
}

// AsError returns nil if the status is a success; otherwise returns an "error" object
// with a concatenated message on reasons of the Status.
func (s *Status) AsError() error {
	if s.IsSuccess() {
		return nil
	}
	if s.err != nil {
		return s.err
	}
	return errors.New(s.Message())
}

// NewStatus makes a Status out of the given arguments and returns its pointer.
func NewStatus(code Code, reasons ...string) *Status {
	s := &Status{
		code:    code,
		reasons: reasons,
	}
	if code == Error {
		s.err = errors.New(s.Message())
	}
	return s
}

// PluginToStatus maps plugin name to status. Currently used to identify which Filter plugin
// returned which status.
type PluginToStatus map[string]*Status

// Merge merges the statuses in the map into one. The resulting status code have the following
// precedence: Error, UnschedulableAndUnresolvable, Unschedulable.
func (p PluginToStatus) Merge() *Status {
	if len(p) == 0 {
		return nil
	}

	finalStatus := NewStatus(Success)
	for _, s := range p {
		if s.Code() == Error {
			finalStatus.err = s.AsError()
		}
		if statusPrecedence[s.Code()] > statusPrecedence[finalStatus.code] {
			finalStatus.code = s.Code()
			// Same as code, we keep the most relevant failedPlugin in the returned Status.
			finalStatus.failedPlugin = s.FailedPlugin()
		}

		for _, r := range s.reasons {
			finalStatus.AppendReason(r)
		}
	}

	return finalStatus
}
*/

// type frameworkImpl struct {
// 	HCPFilterPlugins []HCPFilterPlugin
// 	HCPScorePlugin   []HCPScorePlugin
// }

// func (f *frameworkImpl) RunFilterPlugins(pod *v1.Pod, status *CycleStatus, clusterInfo *resourceinfo.ClusterInfo) *Status {

// }

// func (f *frameworkImpl) runFilterPlugin(pod *v1.Pod, status *CycleStatus, clusterInfo *resourceinfo.ClusterInfo) *Status {
// 	if !status.IsAnyClusters() {
// 		return
// 	}
// 	status := pl.Filter(pod, status, clusterInfo)
// 	return status
// }
