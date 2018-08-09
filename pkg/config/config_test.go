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
	"testing"
)

func TestGetenvReturnsEnvironmentValueIfSet(t *testing.T) {
	// Hopefully PWD is set
	if "nope" == getenv("PWD", "nope") {
		t.Fail()
	}
}

func TestGetenvReturnsDefaultValueIfNotSet(t *testing.T) {
	val := getenv("SMALLPOX", "nobody has smallpox")
	if val != "nobody has smallpox" {
		t.Fail()
	}
}

func defaultsAreReasonable(cfg *Config) bool {
	return (cfg.KubeConfig == "" && cfg.KubeMaster == "" && cfg.MetricsListenAddress == "0.0.0.0:9999" && cfg.Provider == "aws")
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if !defaultsAreReasonable(cfg) {
		t.Fail()
	}
}

// // Calling with no args should get the default config
// func TestParseFlags(t *testing.T) {

// 	KUBE_CONFIG := "/path/to/kube/config"
// 	KUBE_MASTER := "https://master.kubernetes.example.com"
// 	METRICS_ADDR := "127.0.0.1:1234"
// 	REGION := "us-west-2"

// 	args := []string{
// 		fmt.Sprintf("-kube-config %s", KUBE_CONFIG),
// 		fmt.Sprintf("-master-url %s", KUBE_MASTER),
// 		fmt.Sprintf("-metrics-url %s", METRICS_ADDR),
// 		fmt.Sprintf("-region %s", REGION),
// 	}
// 	argv := os.Args
// 	os.Args = args

// 	cfg := DefaultConfig()

// 	if !defaultsAreReasonable(cfg) {
// 		fmt.Println("Defaults are not reasonable")
// 		t.Fail()
// 	}

// 	// No args shouldn't raise an error
// 	if err := cfg.ParseFlags(); err != nil {
// 		fmt.Printf("Error: %s\n", err.Error())
// 		t.Fail()
// 	}

// 	os.Args = argv

// 	fmt.Printf("%v\n", *cfg)
// 	if cfg.KubeConfig != KUBE_CONFIG { t.Fail() }
// 	if cfg.KubeMaster != KUBE_MASTER { t.Fail() }
// 	if cfg.MetricsListenAddress != METRICS_ADDR { t.Fail() }
// 	if cfg.AWSRegion != REGION { t.Fail() }
// }
