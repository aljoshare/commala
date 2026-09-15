package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigParsing_PatternWhitelists(t *testing.T) {
	yamlContent := `
whitelist:
  branches:
    - "^release-please--.*"
    - "^dependabot/.*"
validate:
  branch:
    enabled: true
    whitelist: ["bot@example.com"]
    patterns:
      - "^release-please--.*"
  signoff:
    enabled: true
    whitelist: []
    branch_patterns:
      - "^release-please--.*"
    patterns:
      - '^chore\(.*\): release .*'
  message:
    enabled: true
    whitelist: []
    branch_patterns:
      - "^release-please--.*"
    patterns:
      - '^chore\(.*\): release .*'
`
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, ".commala.yml")
	err := os.WriteFile(configFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	viper.Reset()
	viper.Set("config", configFile)

	c := Config{}
	c.ReadConfig()

	expectedWhitelistBranches := []string{"^release-please--.*", "^dependabot/.*"}
	if !reflect.DeepEqual(c.WhitelistBranches, expectedWhitelistBranches) {
		t.Errorf("expected WhitelistBranches %v, got %v", expectedWhitelistBranches, c.WhitelistBranches)
	}

	expectedBranchPatterns := []string{"^release-please--.*"}
	if !reflect.DeepEqual(c.BranchPatterns, expectedBranchPatterns) {
		t.Errorf("expected BranchPatterns %v, got %v", expectedBranchPatterns, c.BranchPatterns)
	}

	expectedSignOffBranchPatterns := []string{"^release-please--.*"}
	if !reflect.DeepEqual(c.SignOffBranchPatterns, expectedSignOffBranchPatterns) {
		t.Errorf("expected SignOffBranchPatterns %v, got %v", expectedSignOffBranchPatterns, c.SignOffBranchPatterns)
	}

	expectedSignOffPatterns := []string{`^chore\(.*\): release .*`}
	if !reflect.DeepEqual(c.SignOffPatterns, expectedSignOffPatterns) {
		t.Errorf("expected SignOffPatterns %v, got %v", expectedSignOffPatterns, c.SignOffPatterns)
	}

	expectedMessageBranchPatterns := []string{"^release-please--.*"}
	if !reflect.DeepEqual(c.MessageBranchPatterns, expectedMessageBranchPatterns) {
		t.Errorf("expected MessageBranchPatterns %v, got %v", expectedMessageBranchPatterns, c.MessageBranchPatterns)
	}

	expectedMessagePatterns := []string{`^chore\(.*\): release .*`}
	if !reflect.DeepEqual(c.MessagePatterns, expectedMessagePatterns) {
		t.Errorf("expected MessagePatterns %v, got %v", expectedMessagePatterns, c.MessagePatterns)
	}
}
