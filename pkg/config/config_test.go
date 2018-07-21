package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Interval != 30 {
		t.Fail()
	}

	if cfg.KubeConfig != "" || cfg.KubeMaster != "" {
		t.Fail()
	}

	if cfg.Provider != "aws" {
		t.Fail()
	}
}
