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
	"context"
	"time"

	"github.com/cmattoon/aws-ssm/pkg/config"
	"github.com/cmattoon/aws-ssm/pkg/provider"
	"github.com/cmattoon/aws-ssm/pkg/secret"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// Controller is our main struct
type Controller struct {
	Interval      time.Duration
	Provider      provider.Provider
	KubeGen       ClientGenerator
	Context       context.Context
	LabelSelector string
}

// NewController initialises above struct
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
		Context:       context.Background(),
		LabelSelector: cfg.LabelSelector,
	}

	return ctrl
}

// HandleSecrets loops through all k8s api secrets
func (c *Controller) HandleSecrets(cli kubernetes.Interface) error {
	secrets, err := cli.CoreV1().Secrets("").List(c.Context, metav1.ListOptions{LabelSelector: c.LabelSelector})
	if err != nil {
		log.Fatalf("Error retrieving secrets: %s", err)
	}

	i, j, k := 0, 0, 0
	for _, sec := range secrets.Items {
		i++

		obj, err := secret.FromKubernetesSecret(c.Context, c.Provider, sec)
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
		log.Debugf("Successfully updated %s/%s", obj.Namespace, obj.Name)
		k++
	}

	log.Infof("Updated %v/%v secrets (of %v total secrets)", k, j, i)
	return err
}

// WatchSecrets listens for secrets that are created and processes them immediately
func (c *Controller) WatchSecrets(cli kubernetes.Interface) error {
	for {
		select {
		case <-c.Context.Done():
			return nil
		default:
			watcher, err := cli.CoreV1().Secrets(v1.NamespaceAll).Watch(c.Context, metav1.ListOptions{LabelSelector: c.LabelSelector})
			if err != nil {
				log.Errorf("Error retrieving secrets: %s", err)
				return err
			}

			for event := range watcher.ResultChan() {
				sec := event.Object.(*v1.Secret)
				switch event.Type {
				case watch.Added:
					log.Debugf("Secret %s/%s added", sec.ObjectMeta.Namespace, sec.ObjectMeta.Name)
					obj, err := secret.FromKubernetesSecret(c.Context, c.Provider, *sec)
					if err != nil {
						// Error: Irrelevant Secret
						continue
					}

					_, err = obj.UpdateObject(cli)
					if err != nil {
						log.Warnf("Watcher failed to update object %s/%s", obj.Namespace, obj.Name)
						log.Warn(err.Error())
						continue
					}
					log.Infof("Watcher successfully updated %s/%s", obj.Namespace, obj.Name)
				default: // do nothing
				}
			}
		}
	}
}

func (c *Controller) runOnce() error {
	log.Debug("Running...")
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

// Watch listens to secret create API events to create a secret
func (c *Controller) Watch(stopChan <-chan struct{}) {
	cli, err := c.KubeGen.KubeClient()
	if err != nil {
		log.Fatalf("Error with kubernetes client: %s", err)
	}

	log.Info("My Watch begins...")
	err = c.WatchSecrets(cli)
	if err != nil {
		log.Fatalf("Error with WatchSecrets: %s", err)
	}

	for range stopChan {
		log.Info("Ending watch")
		return
	}
}
