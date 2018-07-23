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
