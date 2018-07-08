package collector

import (
	"k8s.io/client-go/kubernetes"
	"github.com/cmattoon/aws-ssm/pkg/config"
)

type Collector struct {
	client kubernetes.Interface
}

func NewCollector(cfg *config.Config) (*Collector){
	c := &Collector{
		client: cfg.
	}
	return c
}
// Searches for annotations
func (c *Collector) GetDefinedSecrets() {
	ingresses, err := c.client.Extensions().Ingresses(sc.namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	ingresses.Items, err = sc.filterByAnnotations(ingresses.Items)
	if err != nil {
		return nil, err
	}
}
