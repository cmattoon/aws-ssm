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
package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/cmattoon/aws-ssm/pkg/config"
	"github.com/cmattoon/aws-ssm/pkg/controller"
)

func main() {
	cfg := config.DefaultConfig()
	if err := cfg.ParseFlags(); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	log.Infof("Using config: %v", cfg)

	stopChan := make(chan struct{}, 1)

	go doMetrics(cfg.MetricsListenAddress)
	go handleSigterm(stopChan)

	ctrl := controller.NewController(cfg)

	ctrl.Run(stopChan)
}

func handleSigterm(stopChan chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Info("Received SIGTERM. Terminating")
	close(stopChan)
}

func doMetrics(address string) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(address, nil))
}
