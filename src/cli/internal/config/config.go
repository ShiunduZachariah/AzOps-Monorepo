package config

import (
	"os"
	"strings"
)

type Config struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
	AuthMode       string
	Output         string
}

func Load() Config {
	return Config{
		SubscriptionID: firstNonEmpty(os.Getenv("AZOPS_SUBSCRIPTION_ID"), os.Getenv("AZURE_SUBSCRIPTION_ID")),
		TenantID:       firstNonEmpty(os.Getenv("AZOPS_TENANT_ID"), os.Getenv("AZURE_TENANT_ID")),
		ClientID:       firstNonEmpty(os.Getenv("AZOPS_CLIENT_ID"), os.Getenv("AZURE_CLIENT_ID")),
		ClientSecret:   firstNonEmpty(os.Getenv("AZOPS_CLIENT_SECRET"), os.Getenv("AZURE_CLIENT_SECRET")),
		AuthMode:       firstNonEmpty(os.Getenv("AZOPS_AUTH_MODE"), "auto"),
		Output:         firstNonEmpty(os.Getenv("AZOPS_OUTPUT"), "plain"),
	}
}

func (c *Config) Normalize() {
	c.SubscriptionID = strings.TrimSpace(c.SubscriptionID)
	c.TenantID = strings.TrimSpace(c.TenantID)
	c.ClientID = strings.TrimSpace(c.ClientID)
	c.ClientSecret = strings.TrimSpace(c.ClientSecret)
	c.AuthMode = strings.ToLower(strings.TrimSpace(c.AuthMode))

	if c.AuthMode == "" {
		c.AuthMode = "auto"
	}
	if c.Output == "" {
		c.Output = "plain"
	}
}

func (c Config) HasServicePrincipalCredentials() bool {
	return c.TenantID != "" && c.ClientID != "" && c.ClientSecret != ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
