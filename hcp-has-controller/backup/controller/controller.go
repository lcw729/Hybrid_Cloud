package controller

import (
	hcphasv1alpha1 "Hybrid_Cloud/pkg/client/resource/v1alpha1/clientset/versioned"
	informer "Hybrid_Cloud/pkg/client/resource/v1alpha1/informers/externalversions/resource/v1alpha1"
	lister "Hybrid_Cloud/pkg/client/resource/v1alpha1/listers/resource/v1alpha1"
	hcphasscheme "Hybrid_Cloud/pkg/client/sync/v1alpha1/clientset/versioned/scheme"

	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	hcphasclientset hcphasv1alpha1.Interface
	hcphasLister    lister.HCPHybridAutoScalerLister
	hcphasSynced    cache.InformerSynced
	workqueue       workqueue.RateLimitingInterface
	recorder        record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	hcphasclientset hcphasv1alpha1.Interface,
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
			utilruntime.HandleError(fmt.Errorf("hcphas '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}
	_ = hcphas
	/*
		// get HCP hybridautoscalers Info
		scaler := autoscaler.NewAutoScaler()
		master_config, _ := cobrautil.BuildConfigFromFlags("master", "/root/.kube/config")
		clientset, _ := hcphasv1alpha1.NewForConfig(master_config)
		warning_count := hcphas.Spec.WarningCount
		options := hcphas.Spec.ScalingOptions
		resource_status := hcphas.Status.ResourceStatus
		target_cluster := hcphas.Spec.TargetCluster
		mode := hcphas.Spec.Mode
		hcpdeployment, _ := clientset.HcpV1alpha1().HCPDeployments("").Get(context.TODO(), name, v1.GetOptions{})
		targets := hcpdeployment.Spec.SchedulingResult.Targets

		if !scaler.ExistDeployment(hcpdeployment) {
			scaler.RegisterDeploymentToAutoScaler(hcpdeployment, hcphas)
		}
		// watching_level 계산
		for _, target := range targets {
			scaler.WarningCountPlusOne(hcpdeployment, target)
			scaler.AutoScaling(hcpdeployment, target)
			fmt.Println("current warningcount is ", scaler.GetWarningCount(hcpdeployment, target))
		}

		// hpa
		hpa_template := options.HpaTemplate
		hpa_namespace := hpa_template.ObjectMeta.Namespace

		// vpa
		vpa_template := options.VpaTemplate
		vpa_namespace := vpa_template.ObjectMeta.Namespace

		// create target_cluster clientset
		config, err := cobrautil.BuildConfigFromFlags(target_cluster, "/root/.kube/config")
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// create vpa_clientset
		vpa_clientset, _ := vpaclientset.NewForConfig(config)
		if mode == "scaling" {
			fmt.Println(mode)
			// check resource_status [ WAITING | DOING | DONE ]
			if resource_status == "WAITING" {

				// check warning_count [ range : 1 <= warning_count <=3 ]
				if warning_count == 1 {

					// create hpa resource for deployment [in target cluster]
					_, err = clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(hpa_namespace).Create(context.TODO(), &hpa_template, metav1.CreateOptions{})
					if err != nil {
						utilruntime.HandleError(fmt.Errorf(err.Error()))
						return err
					} else {
						klog.Info("Success to Create HorizontalPodAutoscalers resource : ", hcphas.ObjectMeta.Name)

						// update has resource status
						hcphas.Status.ResourceStatus = "DONE"
						hcphas.Status.LastSpec = hcphas.Spec
						c.hcphasclientset.HcpV1alpha1().HCPHybridAutoScalers(namespace).Update(context.TODO(), hcphas, metav1.UpdateOptions{})
					}
				} else if warning_count == 2 {

					last_maxReplicas := hcphas.Status.LastSpec.ScalingOptions.HpaTemplate.Spec.MaxReplicas

					// Check if previous MaxReplicas is less than requested value
					if hpa_template.Spec.MaxReplicas > last_maxReplicas {

						// update hpa resource for deployment [in target cluster]
						_, err = clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(hpa_namespace).Update(context.TODO(), &hpa_template, metav1.UpdateOptions{})
						if err != nil {
							utilruntime.HandleError(fmt.Errorf(err.Error()))
							return err
						} else {
							klog.Info("Success to Create HorizontalPodAutoscalers resource : ", hcphas.ObjectMeta.Name)

							// update has resource status
							hcphas.Status.ResourceStatus = "DONE"
							hcphas.Status.LastSpec = hcphas.Spec
							c.hcphasclientset.HcpV1alpha1().HCPHybridAutoScalers(namespace).Update(context.TODO(), hcphas, metav1.UpdateOptions{})
						}
					} else {
						klog.Info("Set a value greater than the current number of replicas")
					}

				} else if warning_count == 3 {

					// create vpa resource for deployment [in target cluster]
					_, err = vpa_clientset.AutoscalingV1beta2().VerticalPodAutoscalers(vpa_namespace).Create(context.TODO(), &vpa_template, metav1.CreateOptions{})
					if err != nil {
						utilruntime.HandleError(fmt.Errorf(err.Error()))
						return err
					} else {
						klog.Info("Success to Create VerticalPodAutoscalers resource : ", hcphas.ObjectMeta.Name)

						// update status
						hcphas.Status.ResourceStatus = "DONE"
						hcphas.Status.LastSpec = hcphas.Spec
						c.hcphasclientset.HcpV1alpha1().HCPHybridAutoScalers(namespace).Update(context.TODO(), hcphas, metav1.UpdateOptions{})
					}
				} else {
					utilruntime.HandleError(fmt.Errorf("warning count is out of range"))
				}
			} else if resource_status == "DOING" {
				utilruntime.HandleError(fmt.Errorf("creating Autoscaler resource : %s", hcphas.ObjectMeta.Name))
			}
		} else if mode == "expanding" {

		}
	*/
	return nil
}
