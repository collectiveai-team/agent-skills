package adapter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeMarkdownSection_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")

	err := mergeMarkdownSection(path, "my-rule-0.1.0", []byte("## My Rule\nDo the thing."))
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.Contains(content, "<!-- BEGIN agent-skills:rule:my-rule-0.1.0 -->") {
		t.Error("missing begin marker")
	}
	if !strings.Contains(content, "## My Rule") {
		t.Error("missing rule content")
	}
	if !strings.Contains(content, "<!-- END agent-skills:rule:my-rule-0.1.0 -->") {
		t.Error("missing end marker")
	}
}

func TestMergeMarkdownSection_ReplaceExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.md")

	// First write
	err := mergeMarkdownSection(path, "my-rule-0.1.0", []byte("## Version 1"))
	if err != nil {
		t.Fatal(err)
	}

	// Second write — should replace, not append
	err = mergeMarkdownSection(path, "my-rule-0.1.0", []byte("## Version 2"))
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if strings.Contains(content, "Version 1") {
		t.Error("old content still present after update")
	}
	if !strings.Contains(content, "Version 2") {
		t.Error("new content not present")
	}
	// Should only have one pair of markers
	if strings.Count(content, "<!-- BEGIN") != 1 {
		t.Error("multiple begin markers found — idempotency broken")
	}
}

func TestMergeMarkdownSection_PreserveOtherContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")

	// Pre-existing content
	os.WriteFile(path, []byte("# My Project Config\n\nSome existing stuff.\n"), 0o644)

	err := mergeMarkdownSection(path, "test-rule", []byte("## Test Rule"))
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.Contains(content, "# My Project Config") {
		t.Error("pre-existing content was lost")
	}
	if !strings.Contains(content, "## Test Rule") {
		t.Error("rule content not added")
	}
}

func TestMergeRuleFile(t *testing.T) {
	dir := t.TempDir()
	rulesDir := filepath.Join(dir, ".cursor", "rules")

	err := mergeRuleFile(rulesDir, "my-rule-0.1.0", []byte("## My Rule\nContent here."))
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(rulesDir, "my-rule-0.1.0.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "## My Rule") {
		t.Error("rule file content missing")
	}
}

func TestMergeJSONMCPServers_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	servers := map[string]interface{}{
		"aws-iac": map[string]interface{}{
			"command": "uvx",
			"args":    []string{"awslabs.mcp-server-iac@latest"},
		},
	}

	err := mergeJSONMCPServers(path, servers)
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	var config map[string]interface{}
	json.Unmarshal(data, &config)

	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		t.Fatal("mcpServers key not found")
	}
	if _, ok := mcpServers["aws-iac"]; !ok {
		t.Error("aws-iac server not found")
	}
}

func TestMergeJSONMCPServers_PreserveExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	// Pre-existing config
	existing := map[string]interface{}{
		"someOtherKey": "value",
		"mcpServers": map[string]interface{}{
			"existing-server": map[string]interface{}{
				"command": "node",
				"args":    []string{"server.js"},
			},
		},
	}
	data, _ := json.MarshalIndent(existing, "", "  ")
	os.WriteFile(path, data, 0o644)

	// Merge new server
	servers := map[string]interface{}{
		"new-server": map[string]interface{}{
			"command": "npx",
			"args":    []string{"@some/server"},
		},
	}

	err := mergeJSONMCPServers(path, servers)
	if err != nil {
		t.Fatal(err)
	}

	data, _ = os.ReadFile(path)
	var config map[string]interface{}
	json.Unmarshal(data, &config)

	if config["someOtherKey"] != "value" {
		t.Error("existing key was lost")
	}

	mcpServers := config["mcpServers"].(map[string]interface{})
	if _, ok := mcpServers["existing-server"]; !ok {
		t.Error("existing server was lost")
	}
	if _, ok := mcpServers["new-server"]; !ok {
		t.Error("new server not added")
	}
}

func TestInterpolateEnv(t *testing.T) {
	// Set a test env var
	os.Setenv("TEST_AGENT_SKILLS_VAR", "from-env")
	defer os.Unsetenv("TEST_AGENT_SKILLS_VAR")

	tests := []struct {
		input    string
		overrides map[string]string
		expected string
	}{
		{
			input:    "${TEST_AGENT_SKILLS_VAR}",
			expected: "from-env",
		},
		{
			input:    "${NONEXISTENT_VAR:-fallback}",
			expected: "fallback",
		},
		{
			input:     "${TEST_AGENT_SKILLS_VAR:-ignored}",
			overrides: map[string]string{"TEST_AGENT_SKILLS_VAR": "override"},
			expected:  "override",
		},
		{
			input:    "prefix-${NONEXISTENT:-default}-suffix",
			expected: "prefix-default-suffix",
		},
	}

	for _, tt := range tests {
		result := InterpolateEnv(tt.input, tt.overrides)
		if result != tt.expected {
			t.Errorf("InterpolateEnv(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
