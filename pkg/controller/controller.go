package controller

import (
	"time"
	
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"github.com/cmattoon/aws-ssm/pkg/config"
	"github.com/cmattoon/aws-ssm/pkg/provider"
	"github.com/cmattoon/aws-ssm/pkg/secret"
)

type Controller struct {
	Interval time.Duration
	Provider provider.Provider
	KubeGen ClientGenerator
}

func NewController(cfg *config.Config) (*Controller) {
	p, err := provider.NewProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to create provider: %s", err)
	}
	
	scg := &SingletonClientGenerator {
		KubeConfig: cfg.KubeConfig,
		KubeMaster: cfg.KubeMaster,
	}
	
	ctrl := &Controller{
		Interval: time.Duration(cfg.Interval) * time.Second,
		Provider: p,
		KubeGen: scg,
	}
	
	return ctrl
}

func (c *Controller) HandleSecrets(cli kubernetes.Interface) (error) {
	secrets, err := cli.CoreV1().Secrets("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error retrieving secrets: %s", err)
	}

	i, j, k := 0, 0, 0	
	for _, sec := range secrets.Items {
		i += 1

		obj, err := secret.FromKubernetesSecret(c.Provider, sec)
		if err != nil {
			// Error: Irrelevant Secret
			continue
		}
		j += 1

		ksec, err := obj.UpdateObject(cli)
		if err != nil {
			log.Warnf("Failed to update object %s/%s", obj.Namespace, obj.Name)
			continue
		}
		log.Infof("Updated %v", ksec)
		k += 1
	}
	
	log.Infof("Updated %v/%v secrets (of %v total secrets)", k, j, i)
	return err
}

func (c *Controller) RunOnce() (error) {
	log.Info("Running...")
	cli, err := c.KubeGen.KubeClient()
	if err != nil {
		log.Fatalf("Error with kubernetes client: %s", err)
	}
	return c.HandleSecrets(cli)
}

func (c *Controller) Run(stopChan <-chan struct{}) {
	ticker := time.NewTicker(c.Interval)
	
	defer ticker.Stop()

	for {
		err := c.RunOnce()
		if err != nil {
			log.Error(err)
		}

		select {
		case <-ticker.C:
		case <-stopChan:
			log.Info("Ending main controller loop")
			return
		}
	}
}
