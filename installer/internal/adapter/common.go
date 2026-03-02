package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// --- Directory copy helpers (ported from resource/copy.go) ---

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// --- Markdown section merging (idempotent rule insertion) ---

func sectionBeginMarker(ruleID string) string {
	return fmt.Sprintf("<!-- BEGIN agent-skills:rule:%s -->", ruleID)
}

func sectionEndMarker(ruleID string) string {
	return fmt.Sprintf("<!-- END agent-skills:rule:%s -->", ruleID)
}

// mergeMarkdownSection inserts or replaces a marked section in a markdown file.
// If the file doesn't exist, it is created. If markers exist, the section is
// replaced; otherwise the section is appended.
func mergeMarkdownSection(filePath, ruleID string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	begin := sectionBeginMarker(ruleID)
	end := sectionEndMarker(ruleID)
	section := begin + "\n" + string(content) + "\n" + end

	existing, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(filePath, []byte(section+"\n"), 0o644)
		}
		return err
	}

	text := string(existing)

	// Try to replace existing section
	pattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(begin) + `.*?` + regexp.QuoteMeta(end))
	if pattern.MatchString(text) {
		text = pattern.ReplaceAllString(text, section)
	} else {
		// Append
		if !strings.HasSuffix(text, "\n") {
			text += "\n"
		}
		text += "\n" + section + "\n"
	}

	return os.WriteFile(filePath, []byte(text), 0o644)
}

// --- Cursor-style rule file merging (one .md file per rule in a directory) ---

// mergeRuleFile writes a rule as an individual markdown file in a directory.
func mergeRuleFile(dirPath, ruleID string, content []byte) error {
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return err
	}
	filePath := filepath.Join(dirPath, ruleID+".md")
	return os.WriteFile(filePath, content, 0o644)
}

// --- JSON config merging (MCP server entries) ---

// mergeJSONMCPServers reads a JSON config file, merges entries into the
// "mcpServers" key, and writes back with indentation. Creates the file if
// it doesn't exist. Preserves all existing keys.
func mergeJSONMCPServers(filePath string, servers map[string]interface{}) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	var config map[string]interface{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			config = make(map[string]interface{})
		} else {
			return err
		}
	} else {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("parsing %s: %w", filePath, err)
		}
	}

	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
	}

	for id, entry := range servers {
		mcpServers[id] = entry
	}
	config["mcpServers"] = mcpServers

	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, append(out, '\n'), 0o644)
}

// --- Env var interpolation ---

var envVarPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?::-(.*?))?\}`)

// InterpolateEnv replaces ${VAR:-default} patterns in a string.
// envOverrides takes precedence, then os.Getenv, then the default value.
func InterpolateEnv(s string, envOverrides map[string]string) string {
	return envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		groups := envVarPattern.FindStringSubmatch(match)
		varName := groups[1]
		defaultVal := groups[2]

		if v, ok := envOverrides[varName]; ok {
			return v
		}
		if v := os.Getenv(varName); v != "" {
			return v
		}
		return defaultVal
	})
}

// buildMCPEntry constructs the JSON-ready map for an MCP server entry.
func buildMCPEntry(command string, args []string, env map[string]string, envOverrides map[string]string) map[string]interface{} {
	entry := map[string]interface{}{
		"command": command,
		"args":    args,
	}
	if len(env) > 0 {
		resolvedEnv := make(map[string]string, len(env))
		for k, v := range env {
			resolvedEnv[k] = InterpolateEnv(v, envOverrides)
		}
		entry["env"] = resolvedEnv
	}
	return entry
}
