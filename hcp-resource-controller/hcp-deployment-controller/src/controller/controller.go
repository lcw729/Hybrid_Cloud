package controller

import (
	resourcev1alpha1apis "Hybrid_Cluster/pkg/apis/resource/v1alpha1"
	resourcev1alpha1 "Hybrid_Cluster/pkg/client/resource/v1alpha1/clientset/versioned"
	resourcev1alpha1scheme "Hybrid_Cluster/pkg/client/resource/v1alpha1/clientset/versioned/scheme"
	informer "Hybrid_Cluster/pkg/client/resource/v1alpha1/informers/externalversions/resource/v1alpha1"
	lister "Hybrid_Cluster/pkg/client/resource/v1alpha1/listers/resource/v1alpha1"
	"context"
	"fmt"
	"time"

	"Hybrid_Cluster/hybridctl/util"

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
	hcpdeploymentclientset resourcev1alpha1.Interface
	hcpdeploymentLister    lister.HCPDeploymentLister
	hcpdeploymentSynced    cache.InformerSynced
	workqueue              workqueue.RateLimitingInterface
	recorder               record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	hcpdeploymentclientset resourcev1alpha1.Interface,
	hcpdeploymentinformer informer.HCPDeploymentInformer) *Controller {
	utilruntime.Must(resourcev1alpha1scheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartStructuredLogging(0)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("hcp")})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:          kubeclientset,
		hcpdeploymentclientset: hcpdeploymentclientset,
		hcpdeploymentLister:    hcpdeploymentinformer.Lister(),
		hcpdeploymentSynced:    hcpdeploymentinformer.Informer().HasSynced,
		workqueue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "hcpdeployment"),
		recorder:               recorder,
	}

	klog.Info("Setting up event handlers")

	hcpdeploymentinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
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
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting HCPDeployment controller")
	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.hcpdeploymentSynced); !ok {
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

	hcpdeployment, err := c.hcpdeploymentLister.HCPDeployments(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("HCPDeployment '%s' in work queue no longer exists", key))
			return nil
		}

		scheduling_status := hcpdeployment.Spec.SchedulingStatus

		if scheduling_status == "Scheduled" {
			targets := hcpdeployment.Spec.SchedulingResult.Targets
			size := len(targets)
			if size < 1 {
				fmt.Println("Target cluster should be more than one.")
			} else {
				var count int32 = 0
				real_replicas := *hcpdeployment.Spec.RealDeploymentSpec.Replicas

				for i := 0; i < size; i++ {
					count += *targets[i].Replicas
				}

				if count < real_replicas {
					fmt.Println("Insufficient number of replicas.")
				} else if count > real_replicas {
					fmt.Println("Excessive number of replicas")
				} else {
					fmt.Println("Appropriate number of replicas")
					ok := DeployKubeDeployment(*hcpdeployment)
					if !ok {
						fmt.Println("fail to schedule deployment")
					} else {
						fmt.Println("success to schedule deployment")

						// HCPDeployment SchedulingStatus 업데이트
						hcpdeployment.Spec.SchedulingStatus = "Completed"
						r, err := c.hcpdeploymentclientset.HcpV1alpha1().HCPDeployments("hcp").Update(context.TODO(), hcpdeployment, metav1.UpdateOptions{})
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Printf("update HCPDeployment %s SchedulingStatus : Completed\n", r.ObjectMeta.Name)
						}
					}
				}
			}
		}

		return err
	}
	return nil
}

func DeployKubeDeployment(hcp_resource resourcev1alpha1apis.HCPDeployment) bool {
	targets := hcp_resource.Spec.SchedulingResult.Targets
	metadata := hcp_resource.Spec.RealDeploymentMetadata
	spec := hcp_resource.Spec.RealDeploymentSpec

	// HCPDeployment SchedulingResult에 따라 Deployment 배포
	for _, target := range targets {
		// cluster clientset 생성

		config, err := util.BuildConfigFromFlags(target.Cluster, "/root/.kube/config")
		if err != nil {
			fmt.Println(err)
			return false
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			return false
		}

		// spec 값 재설정하기
		spec.Replicas = target.Replicas

		// 배포할 Deployment resource 정의
		kube_resource := appsv1.Deployment{
			ObjectMeta: metadata,
			Spec:       spec,
		}

		// Deployment 배포
		r, err := clientset.AppsV1().Deployments(metadata.Namespace).Create(context.TODO(), &kube_resource, metav1.CreateOptions{})

		if err != nil {
			fmt.Println(err)
			return false
		} else {
			fmt.Printf("success to create deployment %s in %s\n", r.ObjectMeta.Name, target.Cluster)
		}
	}
	return true
}
