package controller


import (
	"time"
	
	log "github.com/sirupsen/logrus"
	"github.com/cmattoon/aws-ssm/pkg/config"
	"github.com/cmattoon/aws-ssm/pkg/provider"
)


type Controller struct {
	Interval time.Duration
	Provider provider.Provider
}


func NewController(cfg *config.Config) Controller {
	p, err := provider.NewProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to create provider: %s", err)
	}
	
	ctrl := Controller{
		Interval: time.Duration(cfg.Interval) * time.Second,
		Provider: p,
	}
	
	return ctrl
}

func (c *Controller) RunOnce() error {
	log.Info("Running...")
	name := "com.entic.foo"
	decrypt := true
	
	val, err := c.Provider.GetParameterValue(name, decrypt)
	if err != nil {
		log.Fatalf("Failed to get value: %s", err)
	}
	log.Infof("Got value %s", val)
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
