package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvironmentFilesLoadsInfraEnvAndAppliesAliases(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	tempDir := t.TempDir()
	envDir := filepath.Join(tempDir, "infra", "env")
	if err := os.MkdirAll(envDir, 0o755); err != nil {
		t.Fatalf("create env directory: %v", err)
	}

	envFile := filepath.Join(envDir, ".env")
	envContents := "AZOPS_SUBSCRIPTION_ID=sub-123\nAZOPS_TENANT_ID=tenant-456\nAZOPS_OUTPUT=json\n"
	if err := os.WriteFile(envFile, []byte(envContents), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	t.Setenv("AZOPS_SUBSCRIPTION_ID", "")
	t.Setenv("AZOPS_TENANT_ID", "")
	t.Setenv("AZOPS_OUTPUT", "")
	t.Setenv("AZURE_SUBSCRIPTION_ID", "")
	t.Setenv("AZURE_TENANT_ID", "")

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	if err := LoadEnvironmentFiles(); err != nil {
		t.Fatalf("load environment files: %v", err)
	}

	if got := os.Getenv("AZOPS_SUBSCRIPTION_ID"); got != "sub-123" {
		t.Fatalf("unexpected AZOPS_SUBSCRIPTION_ID: %q", got)
	}

	if got := os.Getenv("AZURE_SUBSCRIPTION_ID"); got != "sub-123" {
		t.Fatalf("unexpected AZURE_SUBSCRIPTION_ID: %q", got)
	}

	if got := os.Getenv("AZURE_TENANT_ID"); got != "tenant-456" {
		t.Fatalf("unexpected AZURE_TENANT_ID: %q", got)
	}
}
