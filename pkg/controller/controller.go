/**
 * Copyright 2018 Curtis Mattoon
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package controller

import (
	"time"

	"github.com/cmattoon/aws-ssm/pkg/config"
	"github.com/cmattoon/aws-ssm/pkg/provider"
	"github.com/cmattoon/aws-ssm/pkg/secret"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	Interval      time.Duration
	Provider      provider.Provider
	KubeGen       ClientGenerator
	LabelSelector string
}

func NewController(cfg *config.Config) *Controller {
	p, err := provider.NewProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to create provider: %s", err)
	}

	scg := &SingletonClientGenerator{
		KubeConfig: cfg.KubeConfig,
		KubeMaster: cfg.KubeMaster,
	}

	ctrl := &Controller{
		Interval:      time.Duration(cfg.Interval) * time.Second,
		Provider:      p,
		KubeGen:       scg,
		LabelSelector: cfg.LabelSelector,
	}

	return ctrl
}

func (c *Controller) HandleSecrets(cli kubernetes.Interface) error {
	secrets, err := cli.CoreV1().Secrets("").List(metav1.ListOptions{LabelSelector: c.LabelSelector})
	if err != nil {
		log.Fatalf("Error retrieving secrets: %s", err)
	}

	i, j, k := 0, 0, 0
	for _, sec := range secrets.Items {
		i++

		obj, err := secret.FromKubernetesSecret(c.Provider, sec)
		if err != nil {
			// Error: Irrelevant Secret
			continue
		}
		j++

		_, err = obj.UpdateObject(cli)
		if err != nil {
			log.Warnf("Failed to update object %s/%s", obj.Namespace, obj.Name)
			log.Warn(err.Error())
			continue
		}
		log.Infof("Successfully updated %s/%s", obj.Namespace, obj.Name)
		k++
	}

	log.Infof("Updated %v/%v secrets (of %v total secrets)", k, j, i)
	return err
}

func (c *Controller) runOnce() error {
	log.Info("Running...")
	cli, err := c.KubeGen.KubeClient()
	if err != nil {
		log.Fatalf("Error with kubernetes client: %s", err)
	}
	return c.HandleSecrets(cli)
}

// Run starts the polling of the k8s API server
func (c *Controller) Run(stopChan <-chan struct{}) {
	ticker := time.NewTicker(c.Interval)

	defer ticker.Stop()

	for {
		err := c.runOnce()
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
