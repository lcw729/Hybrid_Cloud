package controller

import (
	"context"
	"fmt"
	"time"

	resourcev1alpha1clientset "hcp-pkg/client/resource/v1alpha1/clientset/versioned"
	resourcev1alpha1scheme "hcp-pkg/client/resource/v1alpha1/clientset/versioned/scheme"
	Informer "hcp-pkg/client/resource/v1alpha1/informers/externalversions/resource/v1alpha1"
	lister "hcp-pkg/client/resource/v1alpha1/listers/resource/v1alpha1"
	deployment "hcp-pkg/kube-resource/deployment"
	"hcp-pkg/util/clusterManager"

	"hcp-scheduler/src/scheduler"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

const controllerAgentName = "hcp-deployment-controller"

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
	kubeclientset          kubernetes.Interface
	hcpdeploymentclientset resourcev1alpha1clientset.Interface
	hcpdeploymentLister    lister.HCPDeploymentLister
	hcpdeploymentSynced    cache.InformerSynced
	workqueue              workqueue.RateLimitingInterface
	recorder               record.EventRecorder
	scheduler              *scheduler.Scheduler
}

func NewController(
	kubeclientset kubernetes.Interface,
	hcpdeploymentclientset resourcev1alpha1clientset.Interface,
	hcpdeploymentInformer Informer.HCPDeploymentInformer) *Controller {
	utilruntime.Must(resourcev1alpha1scheme.AddToScheme(scheme.Scheme))
	klog.V(4).Infof("Creating event broadcaster")
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartStructuredLogging(0)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("hcp")})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	sched := scheduler.NewScheduler()

	controller := &Controller{
		kubeclientset:          kubeclientset,
		hcpdeploymentclientset: hcpdeploymentclientset,
		hcpdeploymentLister:    hcpdeploymentInformer.Lister(),
		hcpdeploymentSynced:    hcpdeploymentInformer.Informer().HasSynced,
		workqueue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "hcpdeployment"),
		recorder:               recorder,
		scheduler:              sched,
	}

	klog.Infof("Setting up event handlers")

	hcpdeploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueneHCPdeployment,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueneHCPdeployment(new)
		},
	})

	return controller
}

func (c *Controller) enqueneHCPdeployment(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing Informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the Informer factories to begin populating the Informer caches
	klog.Infof("Starting HCPDeployment controller")
	// Wait for the caches to be synced before starting workers
	klog.Infof("Waiting for Informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.hcpdeploymentSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Infof("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Infof("Started workers")
	<-stopCh
	klog.Infof("Shutting down workers")

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
		// workqueue means the items in the Informer cache may actually be
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

	cm, _ := clusterManager.NewClusterManager()
	hcpdeployment, err := c.hcpdeploymentLister.HCPDeployments(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("HCPDeployment '%s' in work queue no longer exists", key))
			return nil
		}
	}

	klog.Infoln("[1] SchedulingNeed", hcpdeployment.Spec.SchedulingNeed)
	klog.Infoln("[2] SchedulingComplete", hcpdeployment.Spec.SchedulingComplete)

	// 스케줄링되지 않은 hcpdeployment 감지
	if !hcpdeployment.Spec.SchedulingNeed && !hcpdeployment.Spec.SchedulingComplete {
		uid, ok := deployment.DeployDeploymentFromHCPDeployment(hcpdeployment)
		if ok {
			klog.Infof("Succeed to deploy deployment %s\n", hcpdeployment.ObjectMeta.Name)
			hcpdeployment.Spec.SchedulingComplete = true
			hcpdeployment.Spec.UUID = uid
			klog.Infof(">>>", uid)
			r, err := c.hcpdeploymentclientset.HcpV1alpha1().HCPDeployments("hcp").Update(context.TODO(), hcpdeployment, metav1.UpdateOptions{})
			if err != nil {
				klog.Error(err)
			} else {
				klog.Infof("Update HCPDeployment %s SchedulingComplete: %t\n", r.ObjectMeta.Name, r.Spec.SchedulingComplete)
			}
		}
	} else if !hcpdeployment.Spec.SchedulingNeed && hcpdeployment.Spec.SchedulingComplete {
		targets := hcpdeployment.Spec.SchedulingResult.Targets
		redeployneed := false
		redeploytarget := map[string]int32{}
		var ns string
		for _, target := range targets {
			if hcpdeployment.Spec.RealDeploymentMetadata.Namespace == "" {
				ns = "default"
			}

			clientset := cm.Cluster_kubeClients[target.Cluster]
			_, err := deployment.GetDeployment(clientset, hcpdeployment.Name, ns)
			if errors.IsNotFound(err) {
				redeployneed = true
				redeploytarget[target.Cluster] = *target.Replicas
			}
		}

		if redeployneed {
			for key, value := range redeploytarget {
				clientset := cm.Cluster_kubeClients[key]
				redeploydeployment := &appsv1.Deployment{}
				redeploydeployment.ObjectMeta = hcpdeployment.Spec.RealDeploymentMetadata
				if redeploydeployment.Namespace == "" {
					redeploydeployment.Namespace = "default"
				}
				redeploydeployment.Spec = hcpdeployment.Spec.RealDeploymentSpec
				redeploydeployment.Spec.Replicas = &value
				err := deployment.CreateDeployment(clientset, "", redeploydeployment)
				if err != nil {
					klog.Error(err)
				} else {
					klog.Infof("Succeed to redeploy deployment %s in %s\n", redeploydeployment.ObjectMeta.Name, key)
					for k := range redeploytarget {
						delete(redeploytarget, k)
					}
				}
			}
		}
	}

	return nil
}
