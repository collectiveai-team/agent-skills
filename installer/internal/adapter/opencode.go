package adapter

import "installer/internal/provider"

func newOpenCode(prov provider.Provider, envOverrides map[string]string) Adapter {
	return &baseAdapter{
		prov:         prov,
		ruleFmt:      ruleStrategyDirectory,
		envOverrides: envOverrides,
	}
}
