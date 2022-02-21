package resource

import (
	cobrautil "Hybrid_Cluster/hybridctl/util"
	hcpclusterv1alpha1 "Hybrid_Cluster/pkg/client/hcpcluster/v1alpha1/clientset/versioned"
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FindNamespaceList(cluster string, namespace string) bool {

	config, err := cobrautil.BuildConfigFromFlags(cluster, "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}

	namespaceList, _ := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	for i := range namespaceList.Items {

		if namespaceList.Items[i].Name == namespace {
			return true
		}
	}
	return false
}

func FindHCPClusterList(cluster string, platform string) bool {
	config, err := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	if err != nil {
		fmt.Println("this error")
	}
	cluster_client := hcpclusterv1alpha1.NewForConfigOrDie(config)

	cluster_list, err := cluster_client.HcpV1alpha1().HCPClusters(platform).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		for _, c := range cluster_list.Items {
			fmt.Println(c.ObjectMeta.Name)
			if c.ObjectMeta.Name == cluster {
				fmt.Printf("find %s in HCPClusterList\n", cluster)
				return true
			}
		}
	}
	fmt.Printf("fail to find %s in HCPClusterList\n", cluster)
	fmt.Println("you should register your cluster to HCP")
	return false
}

func FindHCPClusterList2(cluster string) bool {
	config, err := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	if err != nil {
		fmt.Println("this error")
	}
	cluster_client := hcpclusterv1alpha1.NewForConfigOrDie(config)

	cluster_list, err := cluster_client.HcpV1alpha1().HCPClusters("hcp").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		for _, c := range cluster_list.Items {
			fmt.Println(c.ObjectMeta.Name)
			if c.ObjectMeta.Name == cluster {
				fmt.Printf("find %s in HCPClusterList\n", cluster)
				return true
			}
		}
	}
	fmt.Printf("fail to find %s in HCPClusterList\n", cluster)
	fmt.Println("you should register your cluster to HCP")
	return false
}
