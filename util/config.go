package util

import (
	hcpclusterv1alpha1 "Hybrid_Cluster/pkg/client/hcpcluster/v1alpha1/clientset/versioned"
	"Hybrid_Cluster/util/clusterManager"
	"context"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
			ClientKeyData         string `yaml:"client-key-data,omitempty"`
			Token                 string `yaml:"token,omitempty"`
			AuthProvider          struct {
				Config struct {
					AccessToken string `yaml:"access-token,omitempty"`
					CmdArgs     string `yaml:"cmd-args,omitempty"`
					CmdPath     string `yaml:"cmd-path,omitempty"`
					Expiry      string `yaml:"expiry,omitempty"`
					ExpiryKey   string `yaml:"expiry-key,omitempty"`
					TokenKey    string `yaml:"token-key,omitempty"`
				} `yaml:"config,omitempty"`
				Name string `yaml:"name,omitempty"`
			} `yaml:"auth-provider,omitempty"`
			Exec struct {
				APIVersion string      `yaml:"apiVersion,omitempty"`
				Args       []string    `yaml:"args,omitempty"`
				Command    string      `yaml:"command,omitempty"`
				Env        interface{} `yaml:"env,omitempty"`
			} `yaml:"exec,omitempty"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func UnMarshalKubeConfig(data []byte) (KubeConfig, error) {
	var kubeconfig KubeConfig
	err := yaml.Unmarshal(data, &kubeconfig)
	return kubeconfig, err
}

func ChangeConfigClusterName(platform string, clustername string) error {
	cm := clusterManager.NewClusterManager()
	master_config := cm.Host_config
	hcp_cluster, err := hcpclusterv1alpha1.NewForConfig(master_config)
	if err != nil {
		return err
	}

	cluster, err := hcp_cluster.HcpV1alpha1().HCPClusters(platform).Get(context.TODO(), clustername, metav1.GetOptions{})
	if err != nil {
		return err
	}
	hcpconfig, err := UnMarshalKubeConfig(cluster.Spec.KubeconfigInfo)
	if err != nil {
		return err
	}

	hcpconfig.Clusters[0].Name = clustername
	hcpconfig.Contexts[0].Name = clustername
	hcpconfig.Contexts[0].Context.Cluster = clustername
	hcpconfig.Contexts[0].Context.User = clustername
	hcpconfig.Users[0].Name = clustername

	data, err := yaml.Marshal(hcpconfig)
	if err != nil {
		return err
	}
	cluster.Spec.KubeconfigInfo = data
	_, err = hcp_cluster.HcpV1alpha1().HCPClusters(platform).Update(context.TODO(), cluster, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	// add this config to .kube/config

	bytes, err := ioutil.ReadFile("/root/.kube/config")
	if err != nil {
		return err
	}
	kubeconfig, err := UnMarshalKubeConfig(bytes)
	if err != nil {
		return err
	}
	exist := false
	for _, c := range kubeconfig.Clusters {
		if c.Name == clustername {
			exist = true
			break
		}
	}
	if !exist {
		kubeconfig.Clusters = append(kubeconfig.Clusters, hcpconfig.Clusters...)
		kubeconfig.Contexts = append(kubeconfig.Contexts, hcpconfig.Contexts...)
		kubeconfig.Users = append(kubeconfig.Users, hcpconfig.Users...)
	}
	data, err = yaml.Marshal(&kubeconfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/root/.kube/config", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
