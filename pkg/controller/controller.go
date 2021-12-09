package controller

import (
	sync "Hybrid_Cluster/pkg/apis/sync/v1alpha1"
	v1alpha1hcphas "Hybrid_Cluster/pkg/client/resource/v1alpha1/clientset/versioned"
	informer "Hybrid_Cluster/pkg/client/resource/v1alpha1/informers/externalversions/resource/v1alpha1"
	lister "Hybrid_Cluster/pkg/client/resource/v1alpha1/listers/resource/v1alpha1"
	syncv1alpha1 "Hybrid_Cluster/pkg/client/sync/v1alpha1/clientset/versioned"
	hcphasscheme "Hybrid_Cluster/pkg/client/sync/v1alpha1/clientset/versioned/scheme"
	cm "Hybrid_Cluster/util/clusterManager"
	"context"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const controllerAgentName = "hcp-has-controller"

var syncIndex = 0

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

type Controller struct {
	kubeclientset   kubernetes.Interface
	hcphasclientset v1alpha1hcphas.Interface
	hcphasLister    lister.HCPHybridAutoScalerLister
	hcphasSynced    cache.InformerSynced
	workqueue       workqueue.RateLimitingInterface
	recorder        record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	hcphasclientset v1alpha1hcphas.Interface,
	hcphasinformer informer.HCPHybridAutoScalerInformer) *Controller {
	utilruntime.Must(hcphasscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartStructuredLogging(0)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("hcp")})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:   kubeclientset,
		hcphasclientset: hcphasclientset,
		hcphasLister:    hcphasinformer.Lister(),
		hcphasSynced:    hcphasinformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "hcphas"),
		recorder:        recorder,
	}

	klog.Info("Setting up event handlers")

	hcphasinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueneHCPHAS,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueneHCPHAS(new)
		},
	})

	return controller
}

func (c *Controller) enqueneHCPHAS(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting HCPHybridAutoScaler controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.hcphasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

//

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}
	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	hcphas, err := c.hcphasLister.HCPHybridAutoScalers(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("hcppolicy '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	options := hcphas.Spec.ScalingOptions
	hpa_template := hcphas.Spec.ScalingOptions.HpaTemplate
	if hcphas.Spec.WarningCount == 1 {
		if hcphas.Spec.CurrentStep == "HAS" {
			// Sync 생성
			target_cluster := hcphas.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name
			command := "create"
			h, err := sendSyncHPA(target_cluster, command, options.HpaTemplate)
			if err != nil {
				utilruntime.HandleError(fmt.Errorf(err.Error()))
				return err
			} else {
				klog.Info("Success to Create Sync resource : ", h)
				// hcphas 변경
				hcphas.Status.LastSpec = hcphas.Spec
				hcphas.Spec.CurrentStep = "Sync"
				c.hcphasclientset.HcpV1alpha1().HCPHybridAutoScalers(namespace).Update(context.TODO(), hcphas, metav1.UpdateOptions{})
			}
		}
	} else if hcphas.Spec.WarningCount == 2 {
		if hcphas.Spec.CurrentStep == "HAS" && hpa_template.Spec.MaxReplicas > hcphas.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MaxReplicas {
			// Sync 생성
			target_cluster := hcphas.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name
			command := "update"
			h, err := sendSyncHPA(target_cluster, command, options.HpaTemplate)
			if err != nil {
				utilruntime.HandleError(fmt.Errorf(err.Error()))
				return err
			} else {
				klog.Info("Success to Create Sync resource : ", h)
				// hcphas 변경
				hcphas.Status.LastSpec = hcphas.Spec
				hcphas.Spec.CurrentStep = "Sync"
			}
		}
	} else if hcphas.Spec.WarningCount == 3 {
		if hcphas.Spec.CurrentStep == "VPA" {
			// Sync 생성
			target_cluster := hcphas.Spec.ScalingOptions.HpaTemplate.Spec.ScaleTargetRef.Name
			command := "update"
			h, err := sendSyncHPA(target_cluster, command, options.HpaTemplate)
			if err != nil {
				utilruntime.HandleError(fmt.Errorf(err.Error()))
				return err
			} else {
				klog.Info("Success to Create Sync resource : ", h)
				// hcphas 변경
				hcphas.Status.LastSpec = hcphas.Spec
				hcphas.Spec.CurrentStep = "Sync"
			}
		}
	} else {
		utilruntime.HandleError(fmt.Errorf("warning count is out of range"))
	}
	return nil

}

func sendSyncHPA(clusterName string, command string, template interface{}) (string, error) {
	syncIndex += 1
	cm := cm.NewClusterManager()
	master_config := cm.Host_config
	clientset, err := syncv1alpha1.NewForConfig(master_config)
	if err != nil {
		klog.V(4).Info(err.Error())
	}
	newSync := &sync.Sync{
		ObjectMeta: v1.ObjectMeta{
			Name:      "hcp-hybridautoscaler-hpa-sync-" + strconv.Itoa(syncIndex),
			Namespace: "hcp",
		},
		Spec: sync.SyncSpec{
			ClusterName: clusterName,
			Command:     command,
			Template:    template,
		},
	}
	s, err := clientset.HcpV1alpha1().Syncs("hcp").Create(context.TODO(), newSync, v1.CreateOptions{})
	if err != nil {
		klog.V(4).Info(err.Error())
	}
	klog.V(4).Info("create %s in Namespace %s", s.Name, s.Namespace)
	return s.Name, err
}

// func sendSyncVPA() {

// }
