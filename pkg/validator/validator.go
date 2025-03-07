package validator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Validator validates Kubernetes configurations
type Validator struct {
	ConfigDir string
}

// NewValidator creates a new validator
func NewValidator(configDir string) *Validator {
	return &Validator{
		ConfigDir: configDir,
	}
}

// ValidateAll validates all Kubernetes manifests in the config directory
func (v *Validator) ValidateAll() error {
	// Walk through the config directory
	return filepath.Walk(v.ConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file is a YAML/YML file
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		// Validate the manifest
		if err := v.ValidateFile(path); err != nil {
			return fmt.Errorf("validation failed for %s: %w", path, err)
		}

		return nil
	})
}

// ValidateFile validates a single Kubernetes manifest file
func (v *Validator) ValidateFile(filePath string) error {
	// Use kubectl to validate the file
	cmd := exec.Command("kubectl", "apply", "--dry-run=client", "-f", filePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("kubectl validation failed: %w, output: %s", err, string(output))
	}
	return nil
}

// ValidateHelmChart validates a Helm chart
func (v *Validator) ValidateHelmChart(chartPath string, valuesPath string) error {
	// Use helm to validate the chart
	args := []string{"template", "--debug", chartPath}
	if valuesPath != "" {
		args = append(args, "-f", valuesPath)
	}

	cmd := exec.Command("helm", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("helm validation failed: %w, output: %s", err, string(output))
	}
	return nil
}

// ValidateKustomize validates a Kustomize directory
func (v *Validator) ValidateKustomize(kustomizePath string) error {
	// Use kustomize to validate the directory
	cmd := exec.Command("kustomize", "build", kustomizePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("kustomize validation failed: %w, output: %s", err, string(output))
	}
	return nil
}

// ValidateEnvironment validates all configurations for a specific environment
func (v *Validator) ValidateEnvironment(environment string) error {
	envDir := filepath.Join(v.ConfigDir, environment)

	// Check if the environment directory exists
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return fmt.Errorf("environment directory %s does not exist", envDir)
	}

	v.ConfigDir = envDir
	return v.ValidateAll()
}
