package adapter

import (
	"os"
	"path/filepath"

	"installer/internal/provider"
	"installer/internal/resource"
)

// ruleStrategy controls how rules are written for a provider.
type ruleStrategy int

const (
	// ruleStrategyMergeMarkdown appends rules as marked sections in a single file.
	ruleStrategyMergeMarkdown ruleStrategy = iota
	// ruleStrategyDirectory writes each rule as a separate .md file in a directory.
	ruleStrategyDirectory
)

// defaultMcpKey is the default JSON key used in the config file for MCP servers.
const defaultMcpKey = "mcpServers"

type mcpBuilderFunc func(command string, args []string, env map[string]string, envOverrides map[string]string) map[string]interface{}

// baseAdapter provides the common installation logic shared by all providers.
type baseAdapter struct {
	prov         provider.Provider
	ruleFmt      ruleStrategy
	envOverrides map[string]string
	mcpKey       string
	mcpBuilder   mcpBuilderFunc
}

func (b *baseAdapter) getMcpKey() string {
	if b.mcpKey != "" {
		return b.mcpKey
	}
	return defaultMcpKey
}

func (b *baseAdapter) buildMCP(command string, args []string, env map[string]string) map[string]interface{} {
	if b.mcpBuilder != nil {
		return b.mcpBuilder(command, args, env, b.envOverrides)
	}
	return buildMCPEntry(command, args, env, b.envOverrides)
}

func (b *baseAdapter) ProviderID() string { return b.prov.ID }

func (b *baseAdapter) SupportsResource(resType resource.Type, scope provider.Scope) bool {
	return b.prov.SupportsResource(resType, scope)
}

func (b *baseAdapter) InstallSkill(skill resource.Skill, scope provider.Scope, projectRoot, homeDir string) InstallResult {
	res := InstallResult{Provider: b.prov, Resource: skill, Scope: scope}
	destRoot, ok := b.prov.ResolvePath(resource.TypeSkill, scope, projectRoot, homeDir)
	if !ok {
		res.Skipped = true
		res.Note = "unsupported"
		return res
	}
	dest := filepath.Join(destRoot, skill.Name)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		res.Err = err
		return res
	}
	// Remove existing to ensure idempotency
	_ = os.RemoveAll(dest)
	res.Err = copyDir(skill.Path, dest)
	return res
}

func (b *baseAdapter) InstallRule(rule resource.Rule, scope provider.Scope, projectRoot, homeDir string) InstallResult {
	res := InstallResult{Provider: b.prov, Resource: rule, Scope: scope}
	destPath, ok := b.prov.ResolvePath(resource.TypeRule, scope, projectRoot, homeDir)
	if !ok {
		res.Skipped = true
		res.Note = "unsupported"
		return res
	}

	switch b.ruleFmt {
	case ruleStrategyMergeMarkdown:
		res.Err = mergeMarkdownSection(destPath, rule.Name, rule.Body)
	case ruleStrategyDirectory:
		res.Err = mergeRuleFile(destPath, rule.Name, rule.Body)
	}
	return res
}

func (b *baseAdapter) InstallMCP(mcp resource.MCPServer, scope provider.Scope, projectRoot, homeDir string) InstallResult {
	res := InstallResult{Provider: b.prov, Resource: mcp, Scope: scope}
	destPath, ok := b.prov.ResolvePath(resource.TypeMCPServer, scope, projectRoot, homeDir)
	if !ok {
		res.Skipped = true
		res.Note = "unsupported"
		return res
	}

	servers := make(map[string]interface{}, len(mcp.Servers))
	for _, srv := range mcp.Servers {
		servers[srv.ID] = b.buildMCP(srv.Command, srv.Args, srv.Env)
	}
	res.Err = mergeJSONMCPServers(destPath, b.getMcpKey(), servers)
	return res
}

func (b *baseAdapter) InstallSubagent(sub resource.Subagent, scope provider.Scope, projectRoot, homeDir string) InstallResult {
	res := InstallResult{Provider: b.prov, Resource: sub, Scope: scope}
	destRoot, ok := b.prov.ResolvePath(resource.TypeSubagent, scope, projectRoot, homeDir)
	if !ok {
		res.Skipped = true
		res.Note = "unsupported"
		return res
	}
	dest := filepath.Join(destRoot, sub.DirName)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		res.Err = err
		return res
	}
	_ = os.RemoveAll(dest)
	res.Err = copyDir(sub.Path, dest)
	return res
}
