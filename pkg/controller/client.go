package controller

import (
	"sync"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientGenerator provides clients
type ClientGenerator interface {
        KubeClient() (kubernetes.Interface, error)
}

// SingletonClientGenerator stores provider clients and guarantees that only one instance of client
// will be generated
type SingletonClientGenerator struct {
        KubeConfig string
        KubeMaster string
        client     kubernetes.Interface
        sync.Once
}

// KubeClient generates a kube client if it was not created before
func (p *SingletonClientGenerator) KubeClient() (kubernetes.Interface, error) {
        var err error
        p.Once.Do(func() {
                p.client, err = NewKubeClient(p.KubeConfig, p.KubeMaster)
        })
        return p.client, err
}

// will fallback to restclient.InClusterConfig() if both kubeconfig/master_url == ""
func NewKubeClient(kubeconfig string, master_url string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(master_url, kubeconfig)
	if err != nil {
		return nil, err
	}
	
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	
	log.Infof("Connected to cluster at %s", config.Host)
	return client, nil
}
