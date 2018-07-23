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
	"fmt"
	"testing"

	k8s "k8s.io/client-go/tools/clientcmd"
)

func TestNewKubeClientFailsOnBadFile(t *testing.T) {
	_, err := NewKubeClient("kube-config", "master-url")
	if !(err != nil && err.Error() == "stat kube-config: no such file or directory") {
		t.Fail()
	}
}

func TestNewKubeClientReturnsInClusterConfig(t *testing.T) {
	_, err := NewKubeClient("", "")
	if err.Error() != fmt.Sprintf("invalid configuration: %s", k8s.ErrEmptyConfig.Error()) {
		t.Fail()
	}
}
