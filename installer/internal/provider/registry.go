package provider

import "path/filepath"

func AllProviders() []Provider {
	return []Provider{
		{
			ID:   "claudecode",
			Name: "Claude Code",
			Paths: ResourcePaths{
				SkillsProject:    filepath.Join(".claude", "skills"),
				SkillsGlobal:     filepath.Join(".claude", "skills"),
				RulesProject:     filepath.Join(".claude", "CLAUDE.md"),
				RulesGlobal:      filepath.Join(".claude", "CLAUDE.md"),
				MCPProject:       ".mcp.json",
				MCPGlobal:        ".claude.json",
				SubagentsProject: filepath.Join(".claude", "subagents"),
				SubagentsGlobal:  filepath.Join(".claude", "subagents"),
			},
			Capabilities: Capabilities{
				HasProjectSkills:    true,
				HasGlobalSkills:     true,
				HasProjectRules:     true,
				HasGlobalRules:      true,
				HasProjectMCP:       true,
				HasGlobalMCP:        true,
				HasProjectSubagents: true,
				HasGlobalSubagents:  true,
			},
		},
		{
			ID:   "cursor",
			Name: "Cursor",
			Paths: ResourcePaths{
				SkillsProject:    filepath.Join(".cursor", "skills"),
				SkillsGlobal:     filepath.Join(".cursor", "skills"),
				RulesProject:     filepath.Join(".cursor", "rules"),
				RulesGlobal:      filepath.Join(".cursor", "rules"),
				MCPProject:       filepath.Join(".cursor", "mcp.json"),
				MCPGlobal:        filepath.Join(".cursor", "mcp.json"),
				SubagentsProject: filepath.Join(".cursor", "subagents"),
				SubagentsGlobal:  filepath.Join(".cursor", "subagents"),
			},
			Capabilities: Capabilities{
				HasProjectSkills:    true,
				HasGlobalSkills:     true,
				HasProjectRules:     true,
				HasGlobalRules:      true,
				HasProjectMCP:       true,
				HasGlobalMCP:        true,
				HasProjectSubagents: true,
				HasGlobalSubagents:  true,
			},
		},
		{
			ID:   "windsurf",
			Name: "Windsurf",
			Paths: ResourcePaths{
				SkillsProject:    filepath.Join(".windsurf", "skills"),
				SkillsGlobal:     filepath.Join(".codeium", "windsurf", "skills"),
				RulesProject:     ".windsurfrules",
				RulesGlobal:      filepath.Join(".codeium", "windsurf", "rules"),
				MCPProject:       filepath.Join(".windsurf", "mcp_config.json"),
				MCPGlobal:        filepath.Join(".codeium", "windsurf", "mcp_config.json"),
				SubagentsProject: filepath.Join(".windsurf", "subagents"),
				SubagentsGlobal:  filepath.Join(".codeium", "windsurf", "subagents"),
			},
			Capabilities: Capabilities{
				HasProjectSkills:    true,
				HasGlobalSkills:     true,
				HasProjectRules:     true,
				HasGlobalRules:      true,
				HasProjectMCP:       true,
				HasGlobalMCP:        true,
				HasProjectSubagents: true,
				HasGlobalSubagents:  true,
			},
		},
		{
			ID:   "antigravity",
			Name: "Antigravity",
			Paths: ResourcePaths{
				SkillsProject:    filepath.Join(".agent", "skills"),
				SkillsGlobal:     filepath.Join(".gemini", "antigravity", "skills"),
				RulesProject:     filepath.Join(".agent", "AGENTS.md"),
				RulesGlobal:      filepath.Join(".gemini", "antigravity", "AGENTS.md"),
				MCPProject:       filepath.Join(".agent", "settings.json"),
				MCPGlobal:        filepath.Join(".gemini", "antigravity", "settings.json"),
				SubagentsProject: filepath.Join(".agent", "subagents"),
				SubagentsGlobal:  filepath.Join(".gemini", "antigravity", "subagents"),
			},
			Capabilities: Capabilities{
				HasProjectSkills:    true,
				HasGlobalSkills:     true,
				HasProjectRules:     true,
				HasGlobalRules:      true,
				HasProjectMCP:       true,
				HasGlobalMCP:        true,
				HasProjectSubagents: true,
				HasGlobalSubagents:  true,
			},
		},
		{
			ID:   "gemini",
			Name: "Gemini",
			Paths: ResourcePaths{
				SkillsProject:    filepath.Join(".gemini", "skills"),
				SkillsGlobal:     filepath.Join(".gemini", "skills"),
				RulesProject:     filepath.Join(".gemini", "GEMINI.md"),
				RulesGlobal:      filepath.Join(".gemini", "GEMINI.md"),
				MCPProject:       filepath.Join(".gemini", "settings.json"),
				MCPGlobal:        filepath.Join(".gemini", "settings.json"),
				SubagentsProject: filepath.Join(".gemini", "subagents"),
				SubagentsGlobal:  filepath.Join(".gemini", "subagents"),
			},
			Capabilities: Capabilities{
				HasProjectSkills:    true,
				HasGlobalSkills:     true,
				HasProjectRules:     true,
				HasGlobalRules:      true,
				HasProjectMCP:       true,
				HasGlobalMCP:        true,
				HasProjectSubagents: true,
				HasGlobalSubagents:  true,
			},
		},
		{
			ID:   "opencode",
			Name: "OpenCode",
			Paths: ResourcePaths{
				SkillsProject:    filepath.Join(".opencode", "skills"),
				SkillsGlobal:     filepath.Join(".config", "opencode", "skills"),
				RulesProject:     filepath.Join(".opencode", "rules"),
				RulesGlobal:      filepath.Join(".config", "opencode", "rules"),
				MCPProject:       "opencode.json",
				MCPGlobal:        filepath.Join(".config", "opencode", "opencode.json"),
				SubagentsProject: filepath.Join(".opencode", "subagents"),
				SubagentsGlobal:  filepath.Join(".config", "opencode", "subagents"),
			},
			Capabilities: Capabilities{
				HasProjectSkills:    true,
				HasGlobalSkills:     true,
				HasProjectRules:     true,
				HasGlobalRules:      true,
				HasProjectMCP:       true,
				HasGlobalMCP:        true,
				HasProjectSubagents: true,
				HasGlobalSubagents:  true,
			},
		},
		{
			ID:   "codex",
			Name: "Codex",
			Paths: ResourcePaths{
				SkillsProject:    filepath.Join(".codex", "skills"),
				SkillsGlobal:     filepath.Join(".codex", "skills"),
				RulesProject:     filepath.Join(".codex", "AGENTS.md"),
				RulesGlobal:      filepath.Join(".codex", "AGENTS.md"),
				MCPProject:       filepath.Join(".codex", "config.json"),
				MCPGlobal:        filepath.Join(".codex", "config.json"),
				SubagentsProject: filepath.Join(".codex", "subagents"),
				SubagentsGlobal:  filepath.Join(".codex", "subagents"),
			},
			Capabilities: Capabilities{
				HasProjectSkills:    true,
				HasGlobalSkills:     true,
				HasProjectRules:     true,
				HasGlobalRules:      true,
				HasProjectMCP:       true,
				HasGlobalMCP:        true,
				HasProjectSubagents: true,
				HasGlobalSubagents:  true,
			},
		},
	}
}

// ProviderByID looks up a provider by its ID. Returns nil if not found.
func ProviderByID(id string) *Provider {
	for _, p := range AllProviders() {
		if p.ID == id {
			return &p
		}
	}
	return nil
}
