package provider

import (
	"path/filepath"

	"installer/internal/resource"
)

type Scope int

const (
	ScopeProject Scope = iota
	ScopeGlobal
)

// ResourcePaths defines where each resource type is installed for a provider.
type ResourcePaths struct {
	SkillsProject string
	SkillsGlobal  string
	RulesProject  string
	RulesGlobal   string
	MCPProject    string
	MCPGlobal     string
	SubagentsProject string
	SubagentsGlobal  string
}

// Capabilities declares which resource+scope combos a provider supports.
type Capabilities struct {
	HasProjectSkills    bool
	HasGlobalSkills     bool
	HasProjectRules     bool
	HasGlobalRules      bool
	HasProjectMCP       bool
	HasGlobalMCP        bool
	HasProjectSubagents bool
	HasGlobalSubagents  bool
}

// Provider represents a coding agent/IDE that can receive installed resources.
type Provider struct {
	ID           string
	Name         string
	Paths        ResourcePaths
	Capabilities Capabilities
}

// SupportsResource checks whether this provider supports a given resource type
// at the given scope.
func (p Provider) SupportsResource(resType resource.Type, scope Scope) bool {
	switch resType {
	case resource.TypeSkill:
		if scope == ScopeProject {
			return p.Capabilities.HasProjectSkills
		}
		return p.Capabilities.HasGlobalSkills
	case resource.TypeRule:
		if scope == ScopeProject {
			return p.Capabilities.HasProjectRules
		}
		return p.Capabilities.HasGlobalRules
	case resource.TypeMCPServer:
		if scope == ScopeProject {
			return p.Capabilities.HasProjectMCP
		}
		return p.Capabilities.HasGlobalMCP
	case resource.TypeSubagent:
		if scope == ScopeProject {
			return p.Capabilities.HasProjectSubagents
		}
		return p.Capabilities.HasGlobalSubagents
	}
	return false
}

// ResolvePath returns the absolute destination path for a resource type+scope.
func (p Provider) ResolvePath(resType resource.Type, scope Scope, projectRoot, homeDir string) (string, bool) {
	var rel string
	switch resType {
	case resource.TypeSkill:
		if scope == ScopeProject {
			if p.Capabilities.HasProjectSkills {
				rel = p.Paths.SkillsProject
			}
		} else if p.Capabilities.HasGlobalSkills {
			return filepath.Join(homeDir, p.Paths.SkillsGlobal), true
		}
	case resource.TypeRule:
		if scope == ScopeProject {
			if p.Capabilities.HasProjectRules {
				rel = p.Paths.RulesProject
			}
		} else if p.Capabilities.HasGlobalRules {
			return filepath.Join(homeDir, p.Paths.RulesGlobal), true
		}
	case resource.TypeMCPServer:
		if scope == ScopeProject {
			if p.Capabilities.HasProjectMCP {
				rel = p.Paths.MCPProject
			}
		} else if p.Capabilities.HasGlobalMCP {
			return filepath.Join(homeDir, p.Paths.MCPGlobal), true
		}
	case resource.TypeSubagent:
		if scope == ScopeProject {
			if p.Capabilities.HasProjectSubagents {
				rel = p.Paths.SubagentsProject
			}
		} else if p.Capabilities.HasGlobalSubagents {
			return filepath.Join(homeDir, p.Paths.SubagentsGlobal), true
		}
	}
	if rel == "" {
		return "", false
	}
	return filepath.Join(projectRoot, rel), true
}

// --- Legacy compatibility layer for the existing TUI ---
// These types and functions wrap the new Provider model to maintain the
// interface expected by the current ui/model.go during the transition.

type Agent = Provider

func AllAgents() []Agent {
	return AllProviders()
}

func ResolveDestination(agent Agent, scope Scope, projectRoot string, homeDir string) (string, bool) {
	return agent.ResolvePath(resource.TypeSkill, scope, projectRoot, homeDir)
}
