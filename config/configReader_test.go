package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigLoadsEnvFilesBeforeUnmarshal(t *testing.T) {
	const portEnv = "YOGOURT_CLI_CONFIG_TEST_PORT"

	previousPort, hadPreviousPort := os.LookupEnv(portEnv)
	if err := os.Unsetenv(portEnv); err != nil {
		t.Fatalf("Unsetenv failed: %v", err)
	}
	t.Cleanup(func() {
		if hadPreviousPort {
			_ = os.Setenv(portEnv, previousPort)
			return
		}
		_ = os.Unsetenv(portEnv)
	})

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte(portEnv+"=4242\n"), 0o644); err != nil {
		t.Fatalf("WriteFile .env failed: %v", err)
	}

	configDir := filepath.Join(root, "configs")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll configs failed: %v", err)
	}

	configPath := filepath.Join(configDir, "yogourt.yaml")
	configContent := []byte(`
app_name: test
version: "1.0.0"
mode: development
env_files: .env

server:
  port: ${` + portEnv + `}
  cors: true
`)
	if err := os.WriteFile(configPath, configContent, 0o644); err != nil {
		t.Fatalf("WriteFile config failed: %v", err)
	}

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousDir)
	})

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.Server.Port != 4242 {
		t.Errorf("expected server port 4242, got %d", cfg.Server.Port)
	}
	if len(cfg.EnvFiles) != 1 || cfg.EnvFiles[0] != ".env" {
		t.Errorf("expected env_files to contain .env, got %#v", cfg.EnvFiles)
	}
}
