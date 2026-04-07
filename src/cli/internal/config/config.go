package config

import "os"

type Config struct {
	SubscriptionID string
	Output         string
}

func Load() Config {
	return Config{
		SubscriptionID: firstNonEmpty(os.Getenv("AZOPS_SUBSCRIPTION_ID"), os.Getenv("AZURE_SUBSCRIPTION_ID")),
		Output:         firstNonEmpty(os.Getenv("AZOPS_OUTPUT"), "plain"),
	}
}

func (c *Config) Normalize() {
	if c.Output == "" {
		c.Output = "plain"
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
