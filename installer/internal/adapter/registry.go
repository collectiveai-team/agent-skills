package adapter

import "installer/internal/provider"

type adapterFactory func(provider.Provider, map[string]string) Adapter

var factories = map[string]adapterFactory{
	"claudecode":  newClaudeCode,
	"cursor":      newCursor,
	"windsurf":    newWindsurf,
	"antigravity": newAntigravity,
	"gemini":      newGemini,
	"opencode":    newOpenCode,
	"codex":       newCodex,
}

// ForProvider returns an Adapter for the given provider.
// envOverrides allows callers (e.g. the TUI) to pass user-configured env vars.
func ForProvider(prov provider.Provider, envOverrides map[string]string) Adapter {
	if factory, ok := factories[prov.ID]; ok {
		return factory(prov, envOverrides)
	}
	// Fallback to base adapter for unknown providers
	return &baseAdapter{
		prov:         prov,
		ruleFmt:      ruleStrategyMergeMarkdown,
		envOverrides: envOverrides,
	}
}
