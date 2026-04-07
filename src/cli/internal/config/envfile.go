package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var envAliasMap = map[string]string{
	"AZOPS_TENANT_ID":               "AZURE_TENANT_ID",
	"AZOPS_CLIENT_ID":               "AZURE_CLIENT_ID",
	"AZOPS_CLIENT_SECRET":           "AZURE_CLIENT_SECRET",
	"AZOPS_CLIENT_CERTIFICATE_PATH": "AZURE_CLIENT_CERTIFICATE_PATH",
	"AZOPS_SUBSCRIPTION_ID":         "AZURE_SUBSCRIPTION_ID",
}

func LoadEnvironmentFiles() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working directory: %w", err)
	}

	for _, candidate := range envFileCandidates(workingDir) {
		if err := loadEnvFile(candidate); err != nil {
			return err
		}
	}

	applyEnvAliases()
	return nil
}

func envFileCandidates(startDir string) []string {
	var candidates []string
	seen := map[string]struct{}{}

	for dir := startDir; ; dir = filepath.Dir(dir) {
		addCandidate := func(path string) {
			cleaned := filepath.Clean(path)
			if _, exists := seen[cleaned]; exists {
				return
			}

			seen[cleaned] = struct{}{}
			candidates = append(candidates, cleaned)
		}

		addCandidate(filepath.Join(dir, "infra", "env", ".env"))
		addCandidate(filepath.Join(dir, ".env"))

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	return candidates
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open env file %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		separatorIndex := strings.IndexRune(line, '=')
		if separatorIndex <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:separatorIndex])
		value := strings.TrimSpace(line[separatorIndex+1:])
		if key == "" {
			continue
		}

		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) || (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("set %s from %s line %d: %w", key, path, lineNumber, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read env file %s: %w", path, err)
	}

	return nil
}

func applyEnvAliases() {
	for source, target := range envAliasMap {
		sourceValue := strings.TrimSpace(os.Getenv(source))
		targetValue := strings.TrimSpace(os.Getenv(target))

		if sourceValue != "" && targetValue == "" {
			_ = os.Setenv(target, sourceValue)
		}
	}
}
