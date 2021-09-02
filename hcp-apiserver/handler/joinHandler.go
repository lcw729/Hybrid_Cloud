package handler

import (
	/*
	   clusterRegister "Hybrid_Cluster/clientset/clusterRegister/v1alpha1"
	   "Hybrid_Cluster/hcp-apiserver/converter/mappingTable"
	   "context"
	   "encoding/base64"
	   "flag"
	   "log"
	   "os/exec"
	   "time"

	   "google.golang.org/api/container/v1"
	   "k8s.io/client-go/kubernetes"
	   "k8s.io/client-go/rest"

	   KubeFedCluster "Hybrid_Cluster/apis/kubefedcluster/v1alpha1"

	   "github.com/aws/aws-sdk-go/aws"
	   "github.com/aws/aws-sdk-go/aws/session"
	   corev1 "k8s.io/api/core/v1"


	   "github.com/aws/aws-sdk-go/service/eks"
	   rbacv1 "k8s.io/api/rbac/v1"

	   cobrautil "Hybrid_Cluster/hybridctl/util"

	   "fmt"

	   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	   _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	   "k8s.io/client-go/tools/clientcmd"
	   "k8s.io/client-go/tools/clientcmd/api"
	   "sigs.k8s.io/aws-iam-authenticator/pkg/token"
	   fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	*/
	clusterRegister "Hybrid_Cluster/clientset/v1alpha1"
	"Hybrid_Cluster/hcp-apiserver/converter/mappingTable"
	"context"
	"encoding/base64"
	"flag"
	"log"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"google.golang.org/api/container/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	KubeFedCluster "Hybrid_Cluster/apis/kubefedcluster/v1alpha1"

	corev1 "k8s.io/api/core/v1"

	rbacv1 "k8s.io/api/rbac/v1"

	cobrautil "Hybrid_Cluster/hybridctl/util"

	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

/*
func getAKSClient(authorizer autorest.Authorizer) (containerservice.ManagedClustersClient, error) {
    aksClient := containerservice.NewManagedClustersClient(config.SubscriptionID())
    aksClient.Authorizer = authorizer
    aksClient.AddToUserAgent(config.UserAgent())
    aksClient.PollingDuration = time.Hour * 1
    return aksClient, nil
}

func GetAKS(ctx context.Context, resourceGroupName, resourceName string) (c containerservice.ManagedCluster, err error) {
    aksClient, err := getAKSClient()
    if err != nil {
        return c, fmt.Errorf("cannot get AKS client: %v", err)
    }

    c, err = aksClient.Get(ctx, resourceGroupName, resourceName)
    if err != nil {
        return c, fmt.Errorf("cannot get AKS managed cluster %v from resource group %v: %v", resourceName, resourceGroupName, err)
    }

    return c, nil
}
*/

func Join(info mappingTable.ClusterInfo) bool {

	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	clusterRegisterClientSet, err := clusterRegister.NewForConfig(master_config)
	if err != nil {
		log.Println(err)
	}

	clusterRegisters, err := clusterRegisterClientSet.ClusterRegister(info.PlatformName).Get(info.ClusterName, metav1.GetOptions{})

	if err != nil {
		log.Println(err)
	}

	fmt.Println("---joinHandler start---")
	fmt.Println("--> join process start")

	if info.PlatformName == "gke" {
		/*
		   fmt.Printf("-->  request API to converter [GKE cluster]\n\n")
		   converter.JoinConverter(info)
		*/
		// GKE CLIENTSET
		// var projectId = "keti-container"
		projectId := clusterRegisters.Spec.Projectid
		fProjectId := flag.String("projectId", projectId, "specify a project id to examine")
		flag.Parse()
		if *fProjectId == "" {
			log.Fatal("must specific -projectId")
		}

		// explicit("/home/keti/Downloads/keti-container-7033d28f6fc4.json", projectId)
		kubeConfig, err := getK8sClusterConfigs(context.TODO(), projectId)
		if err != nil {
			log.Println(err)
		}

		var join_cluster_client *kubernetes.Clientset
		var join_cluster_config *rest.Config
		for clusterName := range kubeConfig.Clusters {
			// clusterName "gke_keti-container_us-central1-a_cluster-1"
			gkeClusterName := "gke" + "_" + clusterRegisters.Spec.Projectid + "_" + clusterRegisters.Spec.Region + "_" + info.ClusterName
			if clusterName == gkeClusterName {
				join_cluster_config, err = clientcmd.NewNonInteractiveClientConfig(*kubeConfig, gkeClusterName, &clientcmd.ConfigOverrides{CurrentContext: clusterName}, nil).ClientConfig()
				if err != nil {
					log.Println(err)
				}

				join_cluster_client, err = kubernetes.NewForConfig(join_cluster_config)
				if err != nil {
					log.Println(err)
				}
			}
		}

		JoinCluster(info, join_cluster_client, join_cluster_config.Host)
	} else if info.PlatformName == "aks" {

		cmd := exec.Command("az", "aks", "get-credentials", "--resource-group", clusterRegisters.Spec.Resourcegroup, "--name", clusterRegisters.Spec.Clustername)
		_, err := cmd.Output()
		if err != nil {
			log.Println(err)
		}

		join_cluster_config, _ := cobrautil.BuildConfigFromFlags(info.ClusterName, "/root/.kube/config")
		join_cluster_client := kubernetes.NewForConfigOrDie(join_cluster_config)
		JoinCluster(info, join_cluster_client, join_cluster_config.Host)

	} else if info.PlatformName == "eks" {
		/*
		   fmt.Printf("-->  request API to converter [EKS cluster]\n\n")
		   converter.JoinConverter(info)
		*/
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(clusterRegisters.Spec.Region),
		}))
		eksSvc := eks.New(sess)

		input := &eks.DescribeClusterInput{
			// clusterName eks-master
			Name: aws.String(info.ClusterName),
		}
		result, err := eksSvc.DescribeCluster(input)
		if err != nil {
			fmt.Println(err)
		}

		join_cluster_client, err := newClientset(result.Cluster)
		if err != nil {
			fmt.Println(err)
		}
		JoinCluster(info, join_cluster_client, *result.Cluster.Endpoint)
	}

	/*
	   fmt.Printf("--> connection to VAS\n")
	   httpPostUrl := "http://10.0.5.43:8080/join"
	   reqBody := bytes.NewBufferString("Post plain text")
	   response, err := http.Post(httpPostUrl, "text/plain", reqBody)

	   if err != nil {
	       log.Print(err.Error())
	   }
	   defer response.Body.Close()

	   fmt.Println("response Status:", response.Status)
	   fmt.Println("response Headers:", response.Header)
	   fmt.Println("--> joinCluster func call")
	*/

	fmt.Println("--> join Done!")
	fmt.Println("---joinHandler end---")
	return true
}

/*
kubectl delete ns "kube-federation-system" --context gke_keti-container_us-central1-a_cluster-1;
kubectl delete clusterroles.rbac.authorization.k8s.io kubefed-controller-manager:cluster-1 --context gke_keti-container_us-central1-a_cluster-1 ;
kubectl delete clusterrolebindings.rbac.authorization.k8s.io kubefed-controller-manager:cluster-1-hcp --context gke_keti-container_us-central1-a_cluster-1;
kubectl delete kubefedclusters cluster-1 -n kube-federation-system --context kube-master

kubectl delete ns "kube-federation-system" --context aks-master;
kubectl delete sa aks-master-hcp -n kube-federation-system --context aks-master;
kubectl delete clusterroles.rbac.authorization.k8s.io kubefed-controller-manager:aks-master --context aks-master ;
kubectl delete clusterrolebindings.rbac.authorization.k8s.io kubefed-controller-manager:aks-master-hcp --context aks-master;
kubectl delete kubefedclusters aks-master -n kube-federation-system --context master

kubectl delete ns "kube-federation-system" --context eks-master;
kubectl delete sa eks-master-hcp -n kube-federation-system --context eks-master;
kubectl delete clusterroles.rbac.authorization.k8s.io kubefed-controller-manager:eks-master --context eks-master ;
kubectl delete clusterrolebindings.rbac.authorization.k8s.io kubefed-controller-manager:eks-master-hcp --context eks-master;
kubectl delete kubefedclusters eks-master -n kube-federation-system --context master
*/

func JoinCluster(info mappingTable.ClusterInfo, join_cluster_client *kubernetes.Clientset, APIEndPoint string) bool {

	//var join_cluster_client *kubernetes.Clientset

	master_config, _ := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	// join_cluster_config, _ := cobrautil.BuildConfigFromFlags(info.ClusterName, "/root/.kube/config")

	master_client := kubernetes.NewForConfigOrDie(master_config)
	// join_cluster_client := kubernetes.NewForConfigOrDie(join_cluster_config)

	// 1. CREATE namespace "kube-federation-system"
	Namespace := corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kube-federation-system",
		},
	}

	ns, err_ns := join_cluster_client.CoreV1().Namespaces().Create(context.TODO(), &Namespace, metav1.CreateOptions{})

	if err_ns != nil {
		log.Println(err_ns)
		return false
	} else {
		fmt.Println("< Step 1 > Create Namespace Resource [" + ns.Name + "] in " + info.ClusterName)
	}

	// 2. CREATE service account
	ServiceAccount := corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.ClusterName + "-hcp",
			Namespace: "kube-federation-system",
		},
	}

	sa, err_sa := join_cluster_client.CoreV1().ServiceAccounts("kube-federation-system").Create(context.TODO(), &ServiceAccount, metav1.CreateOptions{})

	if err_sa != nil {
		log.Println(err_sa)
		return false
	} else {
		fmt.Println("< Step 2 > Create Namespace Resource [" + sa.Name + "] in " + info.ClusterName)
	}

	// 3. CREATE cluster role
	ClusterRole := rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kubefed-controller-manager:" + info.ClusterName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{rbacv1.APIGroupAll},
				Verbs:     []string{rbacv1.VerbAll},
				Resources: []string{rbacv1.ResourceAll},
			},
			{
				NonResourceURLs: []string{rbacv1.NonResourceAll},
				Verbs:           []string{"get"},
			},
		},
	}

	cr, err_cr := join_cluster_client.RbacV1().ClusterRoles().Create(context.TODO(), &ClusterRole, metav1.CreateOptions{})

	if err_cr != nil {
		log.Println(err_cr)
		return false
	} else {
		fmt.Println("< Step 3 > Create ClusterRole Resource [" + cr.Name + "] in " + info.ClusterName)
	}

	// 4. CREATE Cluster Role Binding
	ClusterRoleBinding := rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kubefed-controller-manager:" + ServiceAccount.Name,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     ClusterRole.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      ServiceAccount.Name,
				Namespace: ServiceAccount.Namespace,
			},
		},
	}
	_ = ClusterRoleBinding

	crb, err_crb := join_cluster_client.RbacV1().ClusterRoleBindings().Create(context.TODO(), &ClusterRoleBinding, metav1.CreateOptions{})

	if err_crb != nil {
		log.Println(err_crb)
	} else {
		fmt.Println("< Step 4 > Create ClusterRoleBinding Resource [" + crb.Name + "] in " + info.ClusterName)
	}

	time.Sleep(1 * time.Second)

	// 4. GET & CREATE secret (in hcp)
	join_cluster_sa, err_sa1 := join_cluster_client.CoreV1().ServiceAccounts("kube-federation-system").Get(context.TODO(), sa.Name, metav1.GetOptions{})
	if err_sa1 != nil {
		log.Println(err_sa1)
	}
	join_cluster_secret, err_sc := join_cluster_client.CoreV1().Secrets("kube-federation-system").Get(context.TODO(), join_cluster_sa.Secrets[0].Name, metav1.GetOptions{})
	if err_sc != nil {
		log.Println(err_sc)
	} else {
		fmt.Println("< Step 5-1 > Get Secret Resource [" + join_cluster_secret.Name + "] From " + info.ClusterName)
	}

	Secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: info.ClusterName + "-",
			Namespace:    "kube-federation-system",
		},
		Data: map[string][]byte{
			"token": join_cluster_secret.Data["token"],
		},
	}
	cluster_secret, err_secret := master_client.CoreV1().Secrets("kube-federation-system").Create(context.TODO(), Secret, metav1.CreateOptions{})

	if err_secret != nil {
		log.Println(err_secret)
	} else {
		fmt.Println("< Step 5-2 > Create Secret Resource [" + cluster_secret.Name + "] in " + "master")
	}

	// 6. CREATE kubefedcluster (in master)
	kubefedcluster := &fedv1b1.KubeFedCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KubeFedCluster",
			APIVersion: "core.kubefed.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.ClusterName,
			Namespace: "kube-federation-system",
		},
		Spec: fedv1b1.KubeFedClusterSpec{
			// APIEndpoint: *result.Cluster.Endpoint,
			// APIEndpoint: join_cluster_config.Host,
			APIEndpoint: APIEndPoint,
			CABundle:    join_cluster_secret.Data["ca.crt"],
			SecretRef: fedv1b1.LocalSecretReference{
				Name: cluster_secret.Name,
			},
			// DisabledTLSValidations: disabledTLSValidations,
		},
	}

	apiextensionsClientSet, err := KubeFedCluster.NewForConfig(master_config)
	if err != nil {
		log.Println(err)
	}

	newkubefedcluster, err_nkfc := apiextensionsClientSet.KubeFedCluster("kube-federation-system").Create(kubefedcluster)
	if err_nkfc != nil {
		log.Println(err_nkfc)
	} else {
		fmt.Println("< Step 6 > Create KubefedCluster Resource [" + newkubefedcluster.Name + "] in hcp")
	}

	return true
}

func newClientset(cluster *eks.Cluster) (*kubernetes.Clientset, error) {
	opts := &token.GetTokenOptions{
		ClusterID: aws.StringValue(cluster.Name),
	}
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		fmt.Println(err)
	}

	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		fmt.Println(err)
	}

	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		fmt.Println(err)
	}

	clientset, err := kubernetes.NewForConfig(
		&rest.Config{
			Host:        aws.StringValue(cluster.Endpoint),
			BearerToken: tok.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func getK8sClusterConfigs(ctx context.Context, projectId string) (*api.Config, error) {
	svc, err := container.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("container.NewService: %w", err)
	}

	// Basic config structure
	ret := api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters:   map[string]*api.Cluster{},  // Clusters is a map of referencable names to cluster configs
		AuthInfos:  map[string]*api.AuthInfo{}, // AuthInfos is a map of referencable names to user configs
		Contexts:   map[string]*api.Context{},  // Contexts is a map of referencable names to context configs
	}

	// Ask Google for a list of all kube clusters in the given project.
	resp, err := svc.Projects.Zones.Clusters.List(projectId, "-").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("clusters list project=%s: %w", projectId, err)
	}

	for _, f := range resp.Clusters {
		name := fmt.Sprintf("gke_%s_%s_%s", projectId, f.Zone, f.Name)
		cert, err := base64.StdEncoding.DecodeString(f.MasterAuth.ClusterCaCertificate)
		if err != nil {
			return nil, fmt.Errorf("invalid certificate cluster=%s cert=%s: %w", name, f.MasterAuth.ClusterCaCertificate, err)
		}
		// example: gke_my-project_us-central1-b_cluster-1 => https://XX.XX.XX.XX
		ret.Clusters[name] = &api.Cluster{
			CertificateAuthorityData: cert,
			Server:                   "https://" + f.Endpoint,
		}
		// Just reuse the context name as an auth name.
		ret.Contexts[name] = &api.Context{
			Cluster:  name,
			AuthInfo: name,
		}
		// GCP specific configation; use cloud platform scope.
		ret.AuthInfos[name] = &api.AuthInfo{
			AuthProvider: &api.AuthProviderConfig{
				Name: "gcp",
				Config: map[string]string{
					"scopes": "https://www.googleapis.com/auth/cloud-platform",
				},
			},
		}
	}

	return &ret, nil
}

// func explicit(jsonPath, projectID string) {
//     ctx := context.Background()
//     client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
//     if err != nil {
//             log.Fatal(err)
//     }
//     defer client.Close()
//     fmt.Println("Buckets:")
//     it := client.Buckets(ctx, projectID)
//     for {
//             battrs, err := it.Next()
//             if err == iterator.Done {
//                     break
//             }
//             if err != nil {
//                     log.Fatal(err)
//             }
//             fmt.Println(battrs.Name)
//     }
// }
