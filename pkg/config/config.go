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
package config

import (
	"flag"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var Version = "undefined"

func getenv(key string, default_value string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return default_value
	}
	return value
}

type Config struct {
	AWSRegion string
	// Frequency, in seconds, to poll for changes
	Interval             int
	KubeConfig           string
	KubeMaster           string
	MetricsListenAddress string
	Provider             string
	LabelSelector        string
	EnableWatcher        bool
}

func DefaultConfig() *Config {
	cfg := &Config{
		AWSRegion:            "us-west-2",
		Interval:             30,
		KubeConfig:           "",
		KubeMaster:           "",
		MetricsListenAddress: "0.0.0.0:9999",
		Provider:             "aws",
		LabelSelector:        "",
		EnableWatcher:        false,
	}
	return cfg
}

func (cfg *Config) ParseFlags() error {
	kubeConfig := flag.String("kube-config",
		getenv("KUBE_CONFIG", ""),
		"Path to kube config (~/.kube/config)")

	kubeMaster := flag.String("master-url",
		getenv("MASTER_URL", ""),
		"Kubernetes Master URL (kubectl cluster-info)")

	metricAddr := flag.String("metrics-url",
		getenv("METRICS_URL", "0.0.0.0:9999"),
		"Address where metrics/healthz should be served (localhost:9999)")

	region := flag.String("region",
		getenv("AWS_REGION", "us-west-2"),
		"AWS Region (us-west-2)")

	logLevelStr := flag.String("log-level",
		getenv("LOG_LEVEL", "info"),
		"Logrus log level (info)")

	interval := flag.String("interval",
		getenv("SCAN_INTERVAL", "30"),
		"Polling interval")

	labelSelector := flag.String("label-selector",
		getenv("K8S_LABEL_SELECTOR", ""),
		"Label selector for secrets to fetch from k8s API",
	)

	enableWatcher := flag.Bool("watch", false, "Turn on watcher which listens for k8s API events as well as polling")

	flag.Parse()

	i, err := strconv.Atoi(*interval)
	if err != nil {
		log.Error("Could not parse interval - defaulting to 30 seconds")
		i = 30
	}
	// Override config values from CLI
	cfg.AWSRegion = *region
	cfg.Interval = i
	cfg.KubeConfig = *kubeConfig
	cfg.KubeMaster = *kubeMaster
	cfg.MetricsListenAddress = *metricAddr
	cfg.Provider = "aws"
	cfg.LabelSelector = *labelSelector
	cfg.EnableWatcher = *enableWatcher

	logLevel, err := log.ParseLevel(*logLevelStr)
	json := getenv("LOG_FORMAT", "")
	if json == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if err != nil {
		log.Warnf("Improper log level provided: log-level=%s. Defaulting to log-level=info", *logLevelStr)
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	return nil
}
