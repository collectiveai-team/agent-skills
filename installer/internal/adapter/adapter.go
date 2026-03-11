package adapter

import (
	"installer/internal/provider"
	"installer/internal/resource"
)

// InstallResult captures the outcome of a single resource installation.
type InstallResult struct {
	Provider provider.Provider
	Resource resource.Resource
	Scope    provider.Scope
	Err      error
	Skipped  bool
	Note     string
}

// Adapter handles installation of resources for a specific provider.
type Adapter interface {
	ProviderID() string
	InstallSkill(skill resource.Skill, scope provider.Scope, projectRoot, homeDir string) InstallResult
	InstallRule(rule resource.Rule, scope provider.Scope, projectRoot, homeDir string) InstallResult
	InstallMCP(mcp resource.MCPServer, scope provider.Scope, projectRoot, homeDir string) InstallResult
	InstallSubagent(sub resource.Subagent, scope provider.Scope, projectRoot, homeDir string) InstallResult
	SupportsResource(resType resource.Type, scope provider.Scope) bool
}
