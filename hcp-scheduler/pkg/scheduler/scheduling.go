package scheduler

import (
	policy "Hybrid_Cloud/hcp-resource/hcppolicy"
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/priorities"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Scheduler watches for new unscheduled pods. It attempts to find
// nodes that they fit on and writes bindings back to the api server.
type Scheduler struct {
	ClusterClients  map[string]*kubernetes.Clientset
	ClusterInfoList resourceinfo.ClusterInfoList
	SchdPolicy      string
}

func ListTargetClusters() []string {
	cm, _ := clusterManager.NewClusterManager()
	kubefed_clusters := cm.Cluster_list.Items
	var target_clusters []string
	for _, i := range kubefed_clusters {
		target_clusters = append(target_clusters, i.ObjectMeta.Name)
	}
	return target_clusters
}

func createTestNodes() []*NodeInfo {
	// Test
	label1 := map[string]string{"foo": "bar"}
	label2 := map[string]string{"key": "value"}
	label3 := map[string]string{"az": "az1"}
	// label4 := map[string]string{"abc": "az11", "def": "az22"}
	// label5 := map[string]string{"foo": "bar", "key": "value", "az": "az1"}

	affinity1 := &v1.Affinity{
		NodeAffinity: &v1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{{
				Weight: 2,
				Preference: v1.NodeSelectorTerm{
					MatchExpressions: []v1.NodeSelectorRequirement{{
						Key:      "foo",
						Operator: v1.NodeSelectorOpIn,
						Values:   []string{"bar"},
					}},
				},
			}},
		},
	}
	_ = affinity1
	affinity2 := &v1.Affinity{
		NodeAffinity: &v1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
				{
					Weight: 2,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "foo",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"bar"},
							},
						},
					},
				},
				{
					Weight: 4,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "key",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"value"},
							},
						},
					},
				},
				{
					Weight: 5,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "foo",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"bar"},
							},
							{
								Key:      "key",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"value"},
							},
							{
								Key:      "az",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"az1"},
							},
						},
					},
				},
			},
		},
	}

	schd.ClusterInfoList. []*v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "machine1", Labels: label1}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine2", Labels: label2}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine3", Labels: label3}},
	},
}

func NewScheduler() *Scheduler {
	cm, _ := clusterManager.NewClusterManager()

	schd := Scheduler{
		ClusterClients:  cm.Cluster_kubeClients,
		ClusterInfoList: resourceinfo.NewClusterInfoList(),
		// ClusterList:    ListTargetClusters(),
	}

	// HCPPolicy 최적 배치 알고리즘 정책 읽어오기
	algorithm, err := policy.GetAlgorithm()
	if err == nil {
		schd.SchdPolicy = algorithm
	} else {
		schd.SchdPolicy = "DEFAULT_SCHEDPOLICY"
	}

	return &schd
}

type resource struct {
	pod   *v1.Pod
	nodes []*v1.Node
	name  string
}

func (sched *Scheduler) Scheduling(deployment *v1alpha1.HCPDeployment) []v1alpha1.Target {

	test := resource{
		pod: &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{},
			},
			Spec: v1.PodSpec{
				Affinity: affinity2,
			},
		},
		nodes: []*v1.Node{
			{ObjectMeta: metav1.ObjectMeta{Name: "machine1", Labels: label1}},
			{ObjectMeta: metav1.ObjectMeta{Name: "machine2", Labels: label2}},
			{ObjectMeta: metav1.ObjectMeta{Name: "machine3", Labels: label3}},
		},

		name: "all machines are same priority as NodeAffinity is nil",
	}

	schedule_type := sched.SchdPolicy
	replicas := *deployment.Spec.RealDeploymentSpec.Replicas
	var scheduling_result []v1alpha1.Target
	var cnt int32 = 0

	fake_pod := newPodFromHCPDeployment(deployment)

	for i := 0; i < int(replicas); i++ {
		switch schedule_type {
		case "Affinity":
			/*
				target := algorithm.Affinity(&s.ClusterList)
				if target != "" {
					if updateSchedulingResult(&scheduling_result, target) {
						cnt += 1
						fmt.Printf("[Scheduling] %d/%d pod / TargetCluster : %s\n", i+1, replicas, target)
					}
				} else {
					fmt.Println("ERROR: No cluster to be scheduled")
					fmt.Println("Scheduling failed")
					return nil
				}
			*/
			score_table := priorities.NodeAffinity(fake_pod, &sched.ClusterInfoList)
			fmt.Println(score_table.SortScore())
		case "DRF":
		}
	}

	if cnt == replicas {
		fmt.Println("Scheduling succeeded")
		fmt.Println("[Scheduling Result] =====> ", scheduling_result[0].Cluster, *scheduling_result[0].Replicas)
		return scheduling_result
	} else {
		fmt.Println("Scheduling failed")
		return nil
	}
}

func newPodFromHCPDeployment(deployment *v1alpha1.HCPDeployment) *v1.Pod {

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: deployment.Annotations,
			Labels:      deployment.Labels,
		},
		Spec: deployment.Spec.RealDeploymentSpec.Template.Spec,
	}

	return pod
}

func updateSchedulingResult(scheduling_result *[]v1alpha1.Target, cluster string) bool {

	for i, target := range *scheduling_result {
		// 이미 target cluster 목록에 cluster가 있는 경우
		if target.Cluster == cluster {
			// replicas 개수 증가
			temp := *target.Replicas
			temp += 1
			target.Replicas = &temp
			(*scheduling_result)[i] = target
			return true
		}
	}

	// target cluster 목록에 cluster가 없는 경우

	// replicas 개수 1로 설정
	var new_target v1alpha1.Target
	new_target.Cluster = cluster
	var one int32 = 1
	new_target.Replicas = &one
	*scheduling_result = append((*scheduling_result), new_target)
	return true
}

func (sched *Scheduler) scheduleOne(ctx context.Context) {

}
