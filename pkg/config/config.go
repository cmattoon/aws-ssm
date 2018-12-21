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
	//"flag"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	Version = "unknown"
)

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
}

var defaultConfig = &Config{
	AWSRegion:            "us-west-2",
	Interval:             30,
	KubeConfig:           "",
	KubeMaster:           "",
	MetricsListenAddress: "0.0.0.0:9999",
	Provider:             "aws",
}

func NewFromArgs(args []string) *Config {
	cfg := &Config{}
	if err := cfg.ParseFlags(args); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	cfg.Provider = "aws"
	return cfg
}

func (cfg *Config) ParseFlags(args []string) error {
	app := kingpin.New("aws-ssm", "Creates Kubernetes Secrets from AWS SSM Parameter Store")
	app.Version(Version)
	app.DefaultEnvars()

	app.Flag("kube-config", "The kube config file to use").Default(defaultConfig.KubeConfig).StringVar(&cfg.KubeConfig)
	app.Flag("master-url", "The kube master URL from 'kubectl cluster-info'").Default(defaultConfig.KubeMaster).StringVar(&cfg.KubeMaster)
	app.Flag("metrics-url", "The address on which to serve health/metrics ('0.0.0.0:9999')").Default(defaultConfig.MetricsListenAddress).StringVar(&cfg.MetricsListenAddress)
	app.Flag("region", "The AWS region").Default(defaultConfig.AWSRegion).StringVar(&cfg.AWSRegion)
	app.Flag("interval", "The polling interval (in seconds)").Default(string(defaultConfig.Interval)).IntVar(&cfg.Interval)

	_, err := app.Parse(args)
	if err != nil {
		return err
	}
	return nil

}
