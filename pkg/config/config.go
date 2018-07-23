package config

import (
	"os"
	"flag"
)

func getenv(key string, default_value string) (string) {
        value := os.Getenv(key)
        if len(value) == 0 {
                return default_value
        }
        return value
}

type Config struct {
	AWSRegion            string
	// Frequency, in seconds, to poll for changes
	Interval             int
	KubeConfig           string
	KubeMaster           string
	MetricsListenAddress string
	Provider             string
}

func DefaultConfig() *Config {
	cfg := &Config{
		AWSRegion:            "us-west-2",
		Interval:             30,
		KubeConfig:           "",
		KubeMaster:           "",
		MetricsListenAddress: "0.0.0.0:9999",
		Provider:             "aws",
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

        interval := flag.Int("interval", 30, "Polling interval")
        flag.Parse()

	// Override config values from CLI
        cfg.AWSRegion = *region
        cfg.Interval = *interval
        cfg.KubeConfig = *kubeConfig
        cfg.KubeMaster = *kubeMaster
        cfg.MetricsListenAddress = *metricAddr
        cfg.Provider = "aws"

        return nil
}

