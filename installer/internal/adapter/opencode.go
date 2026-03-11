package adapter

import "installer/internal/provider"

func newOpenCode(prov provider.Provider, envOverrides map[string]string) Adapter {
	return &baseAdapter{
		prov:         prov,
		ruleFmt:      ruleStrategyDirectory,
		envOverrides: envOverrides,
		mcpKey:       "mcp",
		mcpBuilder:   buildOpenCodeMCPEntry,
	}
}

func buildOpenCodeMCPEntry(command string, args []string, env map[string]string, envOverrides map[string]string) map[string]interface{} {
	entry := map[string]interface{}{
		"type":    "local",
		"command": append([]string{command}, args...),
	}
	if len(env) > 0 {
		resolvedEnv := make(map[string]string, len(env))
		for k, v := range env {
			resolvedEnv[k] = InterpolateEnv(v, envOverrides)
		}
		entry["environment"] = resolvedEnv
	}
	return entry
}
