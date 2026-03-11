package adapter

import "installer/internal/provider"

func newCursor(prov provider.Provider, envOverrides map[string]string) Adapter {
	return &baseAdapter{
		prov:         prov,
		ruleFmt:      ruleStrategyDirectory,
		envOverrides: envOverrides,
	}
}
