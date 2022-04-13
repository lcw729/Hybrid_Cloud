package scheduler

import (
	policy "Hybrid_Cloud/hcp-resource/hcppolicy"
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

// Scheduler watches for new unscheduled pods. It attempts to find
// nodes that they fit on and writes bindings back to the api server.
type Scheduler struct {
	ClusterClients map[string]*kubernetes.Clientset
	ClusterInfo    []*resourceinfo.ClusterInfo
	ClusterList    []string
	SchdPolicy     string
}

func NewScheduler() *Scheduler {
	cm, _ := clusterManager.NewClusterManager()

	schd := Scheduler{
		ClusterClients: cm.Cluster_kubeClients,
		ClusterInfo:    resourceinfo.NewClusterInfoList(),
	}
	// HCPPolicy 최적 배치 알고리즘 정책 읽어오기
	algorithm, err := policy.GetAlgorithm()
	if err == nil {
		schd.SchdPolicy = algorithm
	} else {
		schd.SchdPolicy = "DEFAULT_SCHEDPOLICY"
	}

	schd.ClusterList = []string{"hcp-cluster", "aks-master"}

	return &schd
}

func (s *Scheduler) Scheduling(deployment *v1alpha1.HCPDeployment) []v1alpha1.Target {
	schedule_type := s.SchdPolicy
	replicas := deployment.Spec.RealDeploymentSpec.Replicas

	fmt.Println(int(*replicas))
	for i := 0; i < int(*replicas); i++ {
		switch schedule_type {
		case "Affinity":
			target := algorithm.Affinity(&s.ClusterList)
			if registerTarget(deployment, target) {
				fmt.Printf("success to scheduler %d/%d pod to %s\n", i+1, int(*replicas), target)
			} else {
				fmt.Println("fail to scheduling")
				return nil
			}
		case "DRF":
		}
	}
	return deployment.Spec.SchedulingResult.Targets
}

func registerTarget(resource *v1alpha1.HCPDeployment, cluster string) bool {

	targets := resource.Spec.SchedulingResult.Targets

	for i, target := range targets {
		// 이미 target cluster 목록에 cluster가 있는 경우
		if target.Cluster == cluster {
			// replicas 개수 증가
			temp := *target.Replicas
			temp += 1
			target.Replicas = &temp
			targets[i] = target
			resource.Spec.SchedulingResult.Targets = targets
			fmt.Println(resource.Spec.SchedulingResult.Targets)
			return true
		}
	}

	// target cluster 목록에 cluster가 없는 경우

	// replicas 개수 1로 설정
	new_target := new(v1alpha1.Target)
	new_target.Cluster = cluster
	var one int32 = 1
	new_target.Replicas = &one
	targets = append(targets, *new_target)
	resource.Spec.SchedulingResult.Targets = targets
	fmt.Println(resource.Spec.SchedulingResult.Targets)

	return true
}

/*
// New returns a Scheduler
func New(client clientset.Interface,
	informerFactory informers.SharedInformerFactory,
	dynInformerFactory dynamicinformer.DynamicSharedInformerFactory,
	recorderFactory profile.RecorderFactory,
	stopCh <-chan struct{},
	opts ...Option) (*Scheduler, error) {

	stopEverything := stopCh
	if stopEverything == nil {
		stopEverything = wait.NeverStop
	}

	options := defaultSchedulerOptions
	for _, opt := range opts {
		opt(&options)
	}

	if options.applyDefaultProfile {
		var versionedCfg v1beta3.KubeSchedulerConfiguration
		scheme.Scheme.Default(&versionedCfg)
		cfg := config.KubeSchedulerConfiguration{}
		if err := scheme.Scheme.Convert(&versionedCfg, &cfg, nil); err != nil {
			return nil, err
		}
		options.profiles = cfg.Profiles
	}

	registry := frameworkplugins.NewInTreeRegistry()
	if err := registry.Merge(options.frameworkOutOfTreeRegistry); err != nil {
		return nil, err
	}

	metrics.Register()

	extenders, err := buildExtenders(options.extenders, options.profiles)
	if err != nil {
		return nil, fmt.Errorf("couldn't build extenders: %w", err)
	}

	podLister := informerFactory.Core().V1().Pods().Lister()
	nodeLister := informerFactory.Core().V1().Nodes().Lister()

	// The nominator will be passed all the way to framework instantiation.
	nominator := internalqueue.NewPodNominator(podLister)
	snapshot := internalcache.NewEmptySnapshot()
	clusterEventMap := make(map[framework.ClusterEvent]sets.String)

	profiles, err := profile.NewMap(options.profiles, registry, recorderFactory,
		frameworkruntime.WithComponentConfigVersion(options.componentConfigVersion),
		frameworkruntime.WithClientSet(client),
		frameworkruntime.WithKubeConfig(options.kubeConfig),
		frameworkruntime.WithInformerFactory(informerFactory),
		frameworkruntime.WithSnapshotSharedLister(snapshot),
		frameworkruntime.WithPodNominator(nominator),
		frameworkruntime.WithCaptureProfile(frameworkruntime.CaptureProfile(options.frameworkCapturer)),
		frameworkruntime.WithClusterEventMap(clusterEventMap),
		frameworkruntime.WithParallelism(int(options.parallelism)),
		frameworkruntime.WithExtenders(extenders),
	)
	if err != nil {
		return nil, fmt.Errorf("initializing profiles: %v", err)
	}

	if len(profiles) == 0 {
		return nil, errors.New("at least one profile is required")
	}

	podQueue := internalqueue.NewSchedulingQueue(
		profiles[options.profiles[0].SchedulerName].QueueSortFunc(),
		informerFactory,
		internalqueue.WithPodInitialBackoffDuration(time.Duration(options.podInitialBackoffSeconds)*time.Second),
		internalqueue.WithPodMaxBackoffDuration(time.Duration(options.podMaxBackoffSeconds)*time.Second),
		internalqueue.WithPodNominator(nominator),
		internalqueue.WithClusterEventMap(clusterEventMap),
		internalqueue.WithPodMaxUnschedulableQDuration(options.podMaxUnschedulableQDuration),
	)

	schedulerCache := internalcache.New(durationToExpireAssumedPod, stopEverything)

	// Setup cache debugger.
	debugger := cachedebugger.New(nodeLister, podLister, schedulerCache, podQueue)
	debugger.ListenForSignal(stopEverything)

	sched := newScheduler(
		schedulerCache,
		extenders,
		internalqueue.MakeNextPodFunc(podQueue),
		MakeDefaultErrorFunc(client, podLister, podQueue, schedulerCache),
		stopEverything,
		podQueue,
		profiles,
		client,
		snapshot,
		options.percentageOfNodesToScore,
	)

	addAllEventHandlers(sched, informerFactory, dynInformerFactory, unionedGVKs(clusterEventMap))

	return sched, nil
}
*/

// scheduleOne does the entire scheduling workflow for a single pod. It is serialized on the scheduling algorithm's host fitting.
func (sched *Scheduler) scheduleOne(ctx context.Context) {

}
