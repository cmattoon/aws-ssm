package controller


import (
	"time"
	
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

func (c *Controller) RunOnce() error {
	log.Info("Running...")
	cli, err := c.KubeGen.KubeClient()
	if err != nil {
		log.Fatalf("Error with kubernetes client: %s", err)
	}

	secrets, err := cli.CoreV1().Secrets("").List(metav1.ListOptions{})
	log.Infof("Found %d secrets\n", len(secrets.Items))
	
	name := "com.entic.foo"
	decrypt := true

	s := secret.NewSecret(c.Provider, name, name, decrypt)
	
	log.Infof("Got value %s", s.Values)
	return nil
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
