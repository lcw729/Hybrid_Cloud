package controller

import (
	cobrautil "Hybrid_Cluster/hybridctl/util"
	sync "Hybrid_Cluster/pkg/apis/sync/v1alpha1"
	v1alpha1Sync "Hybrid_Cluster/pkg/client/sync/v1alpha1/clientset/versioned"
	Syncscheme "Hybrid_Cluster/pkg/client/sync/v1alpha1/clientset/versioned/scheme"
	informers "Hybrid_Cluster/pkg/client/sync/v1alpha1/informers/externalversions/sync/v1alpha1"
	lister "Hybrid_Cluster/pkg/client/sync/v1alpha1/listers/sync/v1alpha1"
	"context"
	"encoding/json"
	"fmt"
	"time"

	vpa "Hybrid_Cluster/pkg/client/vpa/v1beta2/clientset/versioned"

	vpav1beta2 "Hybrid_Cluster/pkg/apis/vpa/v1beta2"

	hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/kubefed/pkg/client/generic"
)

const controllerAgentName = "Sync-controller"

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
	// Kubernetes Core Resource 접근시 사용하는 ClientSet
	kubeclientset kubernetes.Interface
	// Custom Resource 접근시  사용하는 ClientSet
	syncclientset v1alpha1Sync.Interface
	// deploymentsLister  appslisters.DeploymentLister
	// deploymentsSynced  cache.InformerSynced
	syncLister lister.SyncLister
	syncSynced cache.InformerSynced
	workqueue  workqueue.RateLimitingInterface
	recorder   record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	syncclientset v1alpha1Sync.Interface,
	// deploymentInformer appsinformers.DeploymentInformer,
	syncInformer informers.SyncInformer) *Controller {
	utilruntime.Must(Syncscheme.AddToScheme(scheme.Scheme))

	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	// kubernetes client가 클러스터 API를 이용해 내부에 이벤트 전송
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("hcp")})
	// 이벤트 생성
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset: kubeclientset,
		syncclientset: syncclientset,
		// deploymentsLister:  deploymentInformer.Lister(),
		// deploymentsSynced:  deploymentInformer.Informer().HasSynced,
		syncLister: syncInformer.Lister(),
		syncSynced: syncInformer.Informer().HasSynced,
		workqueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Sync"),
		recorder:   recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when resources change
	syncInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueSync,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueSync(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Sync controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.syncSynced); !ok {
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

	s, err := c.syncLister.Syncs(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("sync '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	config, err := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
		return err
	}
	clientset, err := generic.New(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	obj, clusterName, command := c.resourceForSync(s)

	jsonbody, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// if obj.GetKind() == "Deployment" {
	// 	subInstance := &appsv1.Deployment{}
	// 	if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
	// 		// do error check
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	if command == "create" {
	// 		err = clientset.Create(context.TODO(), subInstance)
	// 		if err == nil {
	// 			klog.V(2).Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
	// 		} else {
	// 			klog.V(0).Info("[Error] Cannot Create Deployment : ", err)
	// 		}
	// 	} else if command == "delete" {
	// 		err = clientset.Delete(context.TODO(), subInstance, subInstance.Namespace, subInstance.Name)
	// 		if err == nil {
	// 			klog.V(2).Info("Deleted Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
	// 		} else {
	// 			klog.V(0).Info("[Error] Cannot Delete Deployment : ", err)
	// 		}
	// 	} else if command == "update" {
	// 		err = clientset.Update(context.TODO(), subInstance)
	// 		if err == nil {
	// 			klog.V(2).Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
	// 		} else {
	// 			klog.V(0).Info("[Error] Cannot Update Deployment : ", err)
	// 		}
	// 	}

	// } else
	klog.Info(command)
	klog.Info(obj.GetKind())
	if obj.GetKind() == "HorizontalPodAutoscaler" {
		subInstance := &hpav2beta1.HorizontalPodAutoscaler{}
		if err := json.Unmarshal(jsonbody, &subInstance); err != nil {
			fmt.Println(err)
			return err
		}
		if command == "create" {
			// subInstance.Namespace = "hcp"

			err := clientset.Create(context.TODO(), subInstance)
			if err == nil {
				klog.Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				c.syncclientset.HcpV1alpha1().Syncs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
			} else {
				klog.Info("[Error] Cannot Create HorizontalPodAutoscaler : ", err)
				return err
			}
		} else if command == "update" {
			// subInstance.Namespace = "hcp"
			klog.Info("----")
			err := clientset.Update(context.TODO(), subInstance)
			if err == nil {
				klog.Info("Updated Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + obj.GetNamespace() + "', in Cluster'" + clusterName + "'")
				c.syncclientset.HcpV1alpha1().Syncs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
			} else {
				klog.Info("[Error] Cannot Create HorizontalPodAutoscaler : ", err)
				return err
			}
		}
	} else if obj.GetKind() == "VerticalPodAutoscaler" {
		subInstance := &vpav1beta2.VerticalPodAutoscaler{}
		if err := json.Unmarshal(jsonbody, subInstance); err != nil {
			// do error check
			fmt.Println(err)
			return err
		}
		if command == "create" {
			// subInstance.Namespace = "hcp"
			vpa_clientset, err := vpa.NewForConfig(config)
			vpa, err := vpa_clientset.AutoscalingV1beta2().VerticalPodAutoscalers(subInstance.Namespace).Create(context.TODO(), subInstance, metav1.CreateOptions{})
			if err == nil {
				klog.Info("Created Resource '" + obj.GetKind() + "', Name : '" + obj.GetName() + "',  Namespace : '" + vpa.Namespace + "', in Cluster'" + vpa.ClusterName + "'")
				c.syncclientset.HcpV1alpha1().Syncs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
			} else {
				klog.Info("[Error] Cannot Create VerticalPodAutoscaler : ", err)
				return err
			}
		}
	}
	return nil

}

// enqueueFoo takes a Sync resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Foo.
func (c *Controller) enqueueSync(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

func (c *Controller) resourceForSync(instance *sync.Sync) (*unstructured.Unstructured, string, string) {
	klog.V(4).Info("[Sync] Function Called resourceForSync")
	clusterName := instance.Spec.ClusterName
	command := instance.Spec.Command

	u := &unstructured.Unstructured{}

	klog.V(2).Info("[Parsing Sync] ClusterName : ", clusterName, ", command : ", command)
	var err error
	u.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&instance.Spec.Template)
	if err != nil {
		klog.V(0).Info(err)
	}
	klog.V(4).Info(u.GetName(), " / ", u.GetNamespace())

	return u, clusterName, command
}
