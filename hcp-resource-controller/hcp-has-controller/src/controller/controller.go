package controller

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	resourcev1alpha1 "Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	hcphasv1alpha1 "Hybrid_Cloud/pkg/client/resource/v1alpha1/clientset/versioned"
	informer "Hybrid_Cloud/pkg/client/resource/v1alpha1/informers/externalversions/resource/v1alpha1"
	lister "Hybrid_Cloud/pkg/client/resource/v1alpha1/listers/resource/v1alpha1"
	hcphasscheme "Hybrid_Cloud/pkg/client/sync/v1alpha1/clientset/versioned/scheme"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"
	"time"

	autoscaling "k8s.io/api/autoscaling/v1"
	hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	vpav1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	vpaclientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

var cm, _ = clusterManager.NewClusterManager()

type Controller struct {
	kubeclientset   kubernetes.Interface
	hcphasclientset hcphasv1alpha1.Interface
	hcphasLister    lister.HCPHybridAutoScalerLister
	hcphasSynced    cache.InformerSynced
	workqueue       workqueue.RateLimitingInterface
	recorder        record.EventRecorder
	hasclientset    hcphasv1alpha1.Clientset
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
	hasv1alpha1clientset, _ := hcphasv1alpha1.NewForConfig(cm.Host_config)

	controller := &Controller{
		kubeclientset:   kubeclientset,
		hcphasclientset: hcphasclientset,
		hcphasLister:    hcphasinformer.Lister(),
		hcphasSynced:    hcphasinformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "hcphas"),
		recorder:        recorder,
		hasclientset:    *hasv1alpha1clientset,
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

	// get HCP hybridautoscalers Info
	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clientset, _ := hcphasv1alpha1.NewForConfig(master_config)
	resource_status := hcphas.Status.ResourceStatus
	mode := hcphas.Spec.Mode
	hcpdeployment, _ := clientset.HcpV1alpha1().HCPDeployments("hcp").Get(context.TODO(), name, v1.GetOptions{})
	targets := hcpdeployment.Spec.SchedulingResult.Targets

	if mode == "scaling" {
		fmt.Println(mode)

		var namespace string
		if hcpdeployment.Spec.RealDeploymentMetadata.Namespace == "" {
			namespace = "default"
		} else {
			namespace = hcpdeployment.Spec.RealDeploymentMetadata.Namespace
		}

		var temp *resourcev1alpha1.HCPHybridAutoScaler
		if !hcphas.Status.ScalingInProcess {
			temp = hcphas
			temp.Spec.WarningCount = 0
			temp.Status.ScalingInProcess = true
			temp.Status.ResourceStatus = "WAITING"
			temp.Status.FirstProcess = true
			_, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), temp, v1.UpdateOptions{})
			if err != nil {
				klog.Error(err)
			}
		} else {
			// watching_level 계산
			for _, target := range targets {
				fmt.Println(target.Cluster)
				fmt.Println(hcphas.Status.FirstProcess)
				if resource_status == "DONE" || hcphas.Status.FirstProcess {
					if WatchingLevelCalculator() > 3 {
						hcphas.Spec.WarningCount += 1
						hcphas.Status.FirstProcess = false
						_, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), hcphas, v1.UpdateOptions{})
						if err != nil {
							klog.Error(err)
						}
					}
					switch hcphas.Spec.WarningCount {
					case 1:
						fmt.Printf("===> create new HPA %s in %s\n", hcpdeployment.Name, target.Cluster)
						minReplicas := hcphas.Spec.ScalingOptions.HpaTemplate.Spec.MinReplicas
						maxReplicas := hcphas.Spec.ScalingOptions.HpaTemplate.Spec.MaxReplicas

						// 2. hapTemplate 생성
						hpa := &hpav2beta1.HorizontalPodAutoscaler{
							TypeMeta: v1.TypeMeta{
								Kind:       "HorizontalPodAutoscaler",
								APIVersion: "autoscaling/v2beta1",
							},
							ObjectMeta: v1.ObjectMeta{
								Name:      hcpdeployment.Spec.RealDeploymentMetadata.Name,
								Namespace: namespace,
							},
							Spec: hpav2beta1.HorizontalPodAutoscalerSpec{
								MinReplicas: minReplicas,
								MaxReplicas: maxReplicas,
								ScaleTargetRef: hpav2beta1.CrossVersionObjectReference{
									APIVersion: "apps/v1",
									Kind:       "Deployment",
									Name:       hcpdeployment.Spec.RealDeploymentMetadata.Name,
								},
							},
						}

						// 3. hpaTemplate -> HCPHybridAutoScaler 업데이트
						var temp *resourcev1alpha1.HCPHybridAutoScaler
						temp = hcphas
						temp.Status.LastSpec = hcphas.Spec
						temp.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{HpaTemplate: *hpa}
						temp.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
						newhas, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), temp, v1.UpdateOptions{})
						if err != nil {
							klog.Error(err)
						} else {
							fmt.Println(newhas.Status.ResourceStatus)
							targetclientset := cm.Cluster_kubeClients[target.Cluster]
							newhpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Create(context.TODO(), hpa, v1.CreateOptions{})
							if err != nil {
								klog.Error(err)
							} else {
								klog.Info("Succeed to Create HorizontalPodAutoscalers resource : ", newhpa.ObjectMeta.Name)
								newhas.Status.ResourceStatus = "DONE"
								_, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), newhas, v1.UpdateOptions{})
								if err != nil {
									klog.Error(err)
								}
							}
						}
					case 2:
						fmt.Printf("===> update HPA %s MaxReplicas in %s\n", hcpdeployment.Name, target.Cluster)
						targetclientset := cm.Cluster_kubeClients[target.Cluster]
						hpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(hcpdeployment.Spec.RealDeploymentMetadata.Namespace).Get(context.TODO(), hcpdeployment.Spec.RealDeploymentMetadata.Name, v1.GetOptions{})
						if err != nil {
							klog.Error(err)
						} else {
							var nhpa *hpav2beta1.HorizontalPodAutoscaler = hpa
							nhpa.Spec.MaxReplicas = hpa.Spec.MaxReplicas * 2

							// 3. hpaTemplate -> HCPHybridAutoScaler 업데이트
							hcphas.Status.LastSpec = hcphas.Spec
							hcphas.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{HpaTemplate: *nhpa}
							hcphas.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
							nhas, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), hcphas, v1.UpdateOptions{})
							if err != nil {
								klog.Error(err)
							} else {
								targetclientset := cm.Cluster_kubeClients[target.Cluster]
								newhpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Update(context.TODO(), hpa, v1.UpdateOptions{})
								if err != nil {
									klog.Error(err)
								} else {
									klog.Info("Succeed to Create HorizontalPodAutoscalers resource : ", newhpa.ObjectMeta.Name)
									nhas.Status.ResourceStatus = "DONE"
									_, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), nhas, v1.UpdateOptions{})
									if err != nil {
										klog.Error(err)
									} else {
										fmt.Printf("=====> update %s Done\n", nhas.Name)
									}
								}
							}
						}
					case 3:
						fmt.Printf("===> create new VPA %s in %s\n", hcpdeployment.Name, target.Cluster)

						// 2. vpaTemplate 생성
						updateMode := vpav1beta2.UpdateModeAuto
						vpa := vpav1beta2.VerticalPodAutoscaler{
							TypeMeta: v1.TypeMeta{
								APIVersion: "autoscaling.k8s.io/v1",
								Kind:       "VerticalPodAutoscaler",
							},
							ObjectMeta: v1.ObjectMeta{
								Name:      hcpdeployment.Name,
								Namespace: hcpdeployment.Namespace,
							},
							Spec: vpav1beta2.VerticalPodAutoscalerSpec{
								TargetRef: &autoscaling.CrossVersionObjectReference{
									APIVersion: "apps/v1",
									Kind:       "Deployment",
									Name:       hcpdeployment.Name,
								},
								UpdatePolicy: &vpav1beta2.PodUpdatePolicy{
									UpdateMode: (*vpav1beta2.UpdateMode)(&updateMode),
								},
							},
						}

						// 3. vpaTemplate -> HCPHybridAutoScaler 생성
						hcphas.Status.LastSpec = hcphas.Spec
						hcphas.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{VpaTemplate: vpa}
						hcphas.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
						nhas, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), hcphas, v1.UpdateOptions{})
						if err != nil {
							klog.Error(err)
						} else {
							target_config := cm.Cluster_configs[target.Cluster]
							vpa_clientset, _ := vpaclientset.NewForConfig(target_config)
							_, err := vpa_clientset.AutoscalingV1beta2().VerticalPodAutoscalers(namespace).Create(context.TODO(), &vpa, v1.CreateOptions{})
							if err != nil {
								klog.Error(err)
							} else {
								klog.Info("Success to Create VerticalPodAutoscalers resource : ", hcphas.ObjectMeta.Name)
								nhas.Status.ResourceStatus = "DONE"
								nhas2, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), nhas, v1.UpdateOptions{})
								if err != nil {
									klog.Error(err)
								} else {
									fmt.Printf("=====> update %s Done\n", nhas2.Name)
								}
							}
						}
					}
					// else {
					// 	fmt.Printf("HCPHybridAutoScaler ResourceStatus is not DONE : %s\n", hcphas.Status.ResourceStatus)
					// 	klog.Error("HCPHybridAutoScaler ResourceStatus is not DONE : %s", hcphas.Status.ResourceStatus)
					// }
				}
				fmt.Println("current warningcount is ", hcphas.Spec.WarningCount)
			}
		}
	} else if mode == "expanding" {
		fmt.Println(mode)
	}
	return nil
}

/*
func (c *Controller) CreateHPA(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target, minReplicas *int32, maxReplicas int32) error {
	fmt.Printf("===> create new HPA %s in %s\n", deployment.Name, target.Cluster)

	name := deployment.ObjectMeta.Name
	//	has := a.hasList[name]
	has, _ := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, v1.GetOptions{})
	var namespace string
	if deployment.Spec.RealDeploymentMetadata.Namespace == "" {
		namespace = "default"
	} else {
		namespace = deployment.Spec.RealDeploymentMetadata.Namespace
	}

	// 2. hapTemplate 생성
	hpa := &hpav2beta1.HorizontalPodAutoscaler{
		TypeMeta: v1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      deployment.Spec.RealDeploymentMetadata.Name,
			Namespace: namespace,
		},
		Spec: hpav2beta1.HorizontalPodAutoscalerSpec{
			MinReplicas: minReplicas,
			MaxReplicas: maxReplicas,
			ScaleTargetRef: hpav2beta1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       deployment.Spec.RealDeploymentMetadata.Name,
			},
		},
	}

	// 3. hpaTemplate -> HCPHybridAutoScaler 업데이트
	has.Status.LastSpec = has.Spec
	has.Spec.WarningCount = 1
	has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{HpaTemplate: *hpa}
	has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
	newhas, _ := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, v1.UpdateOptions{})

	targetclientset := cm.Cluster_kubeClients[target.Cluster]
	newhpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Create(context.TODO(), hpa, v1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		klog.Info("Succeed to Create HorizontalPodAutoscalers resource : ", newhpa.ObjectMeta.Name)
		newhas.Status.ResourceStatus = "DONE"
		_, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), newhas, v1.UpdateOptions{})
		if err != nil {
			return err
		} else {
			//a.hasList[name] = has
			return nil
		}
	}
}

func (c *Controller) UpdateHPA(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target) error {
	fmt.Printf("===> update HPA %s MaxReplicas in %s\n", deployment.Name, target.Cluster)
	name := deployment.ObjectMeta.Name
	//	has := a.hasList[name]
	has, _ := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, v1.GetOptions{})
	targetclientset := cm.Cluster_kubeClients[target.Cluster]
	if has.Status.ResourceStatus == "DONE" {
		hpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(deployment.Spec.RealDeploymentMetadata.Namespace).Get(context.TODO(), deployment.Spec.RealDeploymentMetadata.Name, v1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 2-1. hpa max값 설정
		nhpa := hpav2beta1.HorizontalPodAutoscaler{
			TypeMeta: v1.TypeMeta{
				Kind:       "HorizontalPodAutoscaler",
				APIVersion: "autoscaling/v2beta1",
			},
			ObjectMeta: hpa.ObjectMeta,
			Spec: hpav2beta1.HorizontalPodAutoscalerSpec{
				MinReplicas:    hpa.Spec.MinReplicas,
				MaxReplicas:    hpa.Spec.MaxReplicas * 2,
				ScaleTargetRef: hpa.Spec.ScaleTargetRef,
			},
		}
		// 3. hpaTemplate -> HCPHybridAutoScaler 업데이트
		has.Status.LastSpec = has.Spec
		has.Spec.WarningCount = 2
		has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{HpaTemplate: nhpa}
		has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
		//a.hasList[name] = has
		newhas, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, v1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}

		targetclientset := cm.Cluster_kubeClients[target.Cluster]
		newhpa, err := targetclientset.AutoscalingV2beta1().HorizontalPodAutoscalers(deployment.Spec.RealDeploymentMetadata.Namespace).Update(context.TODO(), hpa, v1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			klog.Info("Succeed to Create HorizontalPodAutoscalers resource : ", newhpa.ObjectMeta.Name)
			newhas.Status.ResourceStatus = "DONE"
			_, err := a.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), newhas, v1.UpdateOptions{})
			if err != nil {
				return err
			} else {
				//a.hasList[name] = has
				fmt.Printf("=====> update %s Done\n", newhas.Name)
				return nil
			}
		}
	} else {
		fmt.Printf("HCPHybridAutoScaler ResourceStatus is not DONE : %s\n", has.Status.ResourceStatus)
		return fmt.Errorf("HCPHybridAutoScaler ResourceStatus is not DONE : %s", has.Status.ResourceStatus)
	}
}

func (c *Controller) CreateVPA(deployment resourcev1alpha1.HCPDeployment, target resourcev1alpha1.Target, updateMode string) error {
	fmt.Printf("===> create new VPA %s in %s\n", deployment.Name, target.Cluster)
	name := deployment.ObjectMeta.Name
	has, _ := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, v1.GetOptions{})
	//has := a.hasList[name]
	if has.Status.ResourceStatus == "DONE" {
		// 2. vpaTemplate 생성
		// updateMode := vpav1beta2.UpdateModeAuto
		vpa := vpav1beta2.VerticalPodAutoscaler{
			TypeMeta: v1.TypeMeta{
				APIVersion: "autoscaling.k8s.io/v1",
				Kind:       "VerticalPodAutoscaler",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      deployment.Name,
				Namespace: deployment.Namespace,
			},
			Spec: vpav1beta2.VerticalPodAutoscalerSpec{
				TargetRef: &autoscaling.CrossVersionObjectReference{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       deployment.Name,
				},
				UpdatePolicy: &vpav1beta2.PodUpdatePolicy{
					UpdateMode: (*vpav1beta2.UpdateMode)(&updateMode),
				},
			},
		}

		has, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 3. vpaTemplate -> HCPHybridAutoScaler 생성
		has.Status.LastSpec = has.Spec
		has.Spec.WarningCount = 3
		has.Spec.ScalingOptions = resourcev1alpha1.ScalingOptions{VpaTemplate: vpa}
		has.Status = resourcev1alpha1.HCPHybridAutoScalerStatus{ResourceStatus: "WAITING"}
		//	a.hasList[name] = *has
		newhas, _ := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), has, v1.UpdateOptions{})

		target_config := cm.Cluster_configs[target.Cluster]
		vpa_clientset, _ := vpaclientset.NewForConfig(target_config)
		hcphas, err := vpa_clientset.AutoscalingV1beta2().VerticalPodAutoscalers(deployment.Spec.RealDeploymentMetadata.Namespace).Create(context.TODO(), &vpa, v1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			klog.Info("Success to Create VerticalPodAutoscalers resource : ", hcphas.ObjectMeta.Name)
			newhas.Status.ResourceStatus = "DONE"
			_, err := c.hasclientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Update(context.TODO(), newhas, v1.UpdateOptions{})
			if err != nil {
				return err
			} else {
				//a.hasList[name] = *has
				fmt.Printf("=====> update %s Done\n", newhas.Name)
				return nil
			}
		}
	} else {
		fmt.Printf("HCPHybridAutoScaler ResourceStatus is not DONE : %s\n", has.Status.ResourceStatus)
		return fmt.Errorf("HCPHybridAutoScaler ResourceStatus is not DONE : %s", has.Status.ResourceStatus)
	}
}
*/

func WatchingLevelCalculator() int {
	time.Sleep(time.Second * 10)
	return 4
}
