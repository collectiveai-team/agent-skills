package agents

import "path/filepath"

type Scope int

const (
	ScopeProject Scope = iota
	ScopeGlobal
)

type Agent struct {
	ID          string
	Name        string
	ProjectPath string
	GlobalPath  string
	HasProject  bool
	HasGlobal   bool
	Description string
}

func AllAgents() []Agent {
	return []Agent{
		{
			ID:          "windsurf",
			Name:        "Windsurf",
			ProjectPath: filepath.Join(".windsurf", "skills"),
			GlobalPath:  filepath.Join(".codeium", "windsurf", "skills"),
			HasProject:  true,
			HasGlobal:   true,
		},
		{
			ID:          "cursor",
			Name:        "Cursor",
			ProjectPath: filepath.Join(".cursor", "skills"),
			GlobalPath:  filepath.Join(".cursor", "skills"),
			HasProject:  true,
			HasGlobal:   true,
		},
		{
			ID:          "claudecode",
			Name:        "Claude Code",
			ProjectPath: filepath.Join(".claude", "skills"),
			GlobalPath:  filepath.Join(".claude", "skills"),
			HasProject:  true,
			HasGlobal:   true,
		},
		{
			ID:          "antigravity",
			Name:        "Antigravity",
			ProjectPath: filepath.Join(".agent", "skills"),
			GlobalPath:  filepath.Join(".gemini", "antigravity", "skills"),
			HasProject:  true,
			HasGlobal:   true,
		},
		{
			ID:         "opencode",
			Name:       "OpenCode",
			GlobalPath: filepath.Join(".config", "opencode", "skills"),
			HasProject: false,
			HasGlobal:  true,
		},
		{
			ID:          "gemini",
			Name:        "Gemini",
			ProjectPath: filepath.Join(".gemini", "skills"),
			GlobalPath:  filepath.Join(".gemini", "skills"),
			HasProject:  true,
			HasGlobal:   true,
		},
		{
			ID:          "codex",
			Name:        "Codex",
			ProjectPath: filepath.Join(".codex", "skills"),
			GlobalPath:  filepath.Join(".codex", "skills"),
			HasProject:  true,
			HasGlobal:   true,
		},
	}
}

func ResolveDestination(agent Agent, scope Scope, projectRoot string, homeDir string) (string, bool) {
	if scope == ScopeProject {
		if agent.HasProject {
			return filepath.Join(projectRoot, agent.ProjectPath), true
		}
		if agent.HasGlobal {
			return filepath.Join(homeDir, agent.GlobalPath), true
		}
		return "", false
	}

	if agent.HasGlobal {
		return filepath.Join(homeDir, agent.GlobalPath), true
	}

	if agent.HasProject {
		return filepath.Join(projectRoot, agent.ProjectPath), true
	}

	return "", false
}
