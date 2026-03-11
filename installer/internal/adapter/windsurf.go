package adapter

import "installer/internal/provider"

func newWindsurf(prov provider.Provider, envOverrides map[string]string) Adapter {
	return &baseAdapter{
		prov:         prov,
		ruleFmt:      ruleStrategyMergeMarkdown,
		envOverrides: envOverrides,
	}
}
