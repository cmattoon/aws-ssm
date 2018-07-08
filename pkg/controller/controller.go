package controller


import (
	"time"
	
	log "github.com/sirupsen/logrus"
	"github.com/cmattoon/aws-param-store/pkg/config"
)


type Controller struct {
	Interval time.Duration
}

func NewController(cfg *config.Config) Controller {
	
	ctrl := Controller{
		Interval: time.Duration(cfg.Interval) * time.Second,
	}
	return ctrl
}

func (c *Controller) RunOnce() error {
	log.Info("Running...")
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
