package config

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
		KubeConfig:           "/Users/cmattoon/.kube/config",
		KubeMaster:           "https://api.dev-apps.us-west-2.k8s.entic-int.com",
		MetricsListenAddress: "127.0.0.1:9999",
		Provider:             "aws",
	}
	return cfg
}

