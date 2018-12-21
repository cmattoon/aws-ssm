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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	overriddenConfig = &Config{
		AWSRegion:            "eu-central-1",
		Interval:             60,
		KubeConfig:           "/opt/kube.config",
		KubeMaster:           "http://master.k8s.local",
		MetricsListenAddress: "0.0.0.0:1234",
		Provider:             "aws",
	}
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

func TestFlagsOverrideConfig(t *testing.T) {
	for _, tc := range []struct {
		title    string
		env      map[string]string
		args     []string
		expected *Config
	}{
		{
			title: "override everything with args (kingpin)",
			env:   map[string]string{},
			args: []string{
				"--kube-config=/opt/kube.config",
				"--master-url=http://master.k8s.local",
				"--metrics-url=0.0.0.0:1234",
				"--region=eu-central-1",
				"--interval=60",
			},
			expected: overriddenConfig,
		},
		{
			title: "override everything with env (kingpin)",
			env: map[string]string{
				"AWS_SSM_KUBE_CONFIG": "/opt/kube.config",
				"AWS_SSM_MASTER_URL":  "http://master.k8s.local",
				"AWS_SSM_METRICS_URL": "0.0.0.0:1234",
				"AWS_SSM_REGION":      "eu-central-1",
				"AWS_SSM_INTERVAL":    "60",
			},
			args:     []string{},
			expected: overriddenConfig,
		},
		// {
		// 	title: "override everything with args (deprecated)",
		// 	env:   map[string]string{},
		// 	args: []string{
		// 		"-kube-config=/opt/kube.config",
		// 		"-master-url=http://master.k8s.local",
		// 		"-metrics-url=0.0.0.0:1234",
		// 		"-region=eu-central-1",
		// 		"-interval=60",
		// 	},
		// 	expected: overriddenConfig,
		// },
		// {
		// 	title: "override everything with env (deprecated)",
		// 	env: map[string]string{
		// 		"KUBE_CONFIG":      "/opt/kube.config",
		// 		"MASTER_URL":       "http://master.k8s.local",
		// 		"METRICS_URL":      "0.0.0.0:1234",
		// 		"AWS_REGION":       "eu-central-1",
		// 		"AWS_SSM_INTERVAL": "60", // This couldn't be set before
		// 	},
		// 	args:     []string{},
		// 	expected: overriddenConfig,
		// },
	} {
		t.Run(tc.title, func(t *testing.T) {
			env0 := _setenv(t, tc.env)
			defer func() { _restore_env(t, env0) }()

			cfg := NewFromArgs(tc.args)
			assert.Equal(t, tc.expected, cfg)
		})
	}
}

func _setenv(t *testing.T, e0 map[string]string) map[string]string {
	e1 := map[string]string{}

	for k, v := range e0 {
		e1[k] = os.Getenv(v)
		require.NoError(t, os.Setenv(k, v))
	}
	return e1
}

func _restore_env(t *testing.T, e map[string]string) {
	for k, v := range e {
		require.NoError(t, os.Setenv(k, v))
	}
}
