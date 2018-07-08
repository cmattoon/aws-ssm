package config

type Config struct {
	// Frequency, in seconds, to poll for changes
	Interval             int
	KubeConfig           string
	KubeMaster           string
	MetricsListenAddress string
}

func DefaultConfig() *Config {
	cfg := &Config{
		Interval:             30,
		KubeConfig:           "",
		KubeMaster:           "",
		MetricsListenAddress: "127.0.0.1:9999",
	}
	return cfg
}

