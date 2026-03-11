package installer

import (
	"installer/internal/adapter"
	"installer/internal/provider"
	"installer/internal/resource"
)

// InstallTask describes a single resource to install for a single provider.
type InstallTask struct {
	Provider provider.Provider
	Resource resource.Resource
	Scope    provider.Scope
}

// InstallPlan holds an ordered list of install tasks.
type InstallPlan struct {
	Tasks []InstallTask
}

// BuildPlan creates an install plan for the given providers and resources.
// Unsupported resource+provider+scope combinations are silently skipped.
func BuildPlan(
	providers []provider.Provider,
	skills []resource.Skill,
	rules []resource.Rule,
	mcps []resource.MCPServer,
	subagents []resource.Subagent,
	scope provider.Scope,
) *InstallPlan {
	plan := &InstallPlan{}

	for _, prov := range providers {
		for _, s := range skills {
			if prov.SupportsResource(resource.TypeSkill, scope) {
				plan.Tasks = append(plan.Tasks, InstallTask{
					Provider: prov,
					Resource: s,
					Scope:    scope,
				})
			}
		}
		for _, r := range rules {
			if prov.SupportsResource(resource.TypeRule, scope) {
				plan.Tasks = append(plan.Tasks, InstallTask{
					Provider: prov,
					Resource: r,
					Scope:    scope,
				})
			}
		}
		for _, m := range mcps {
			if prov.SupportsResource(resource.TypeMCPServer, scope) {
				plan.Tasks = append(plan.Tasks, InstallTask{
					Provider: prov,
					Resource: m,
					Scope:    scope,
				})
			}
		}
		for _, s := range subagents {
			if prov.SupportsResource(resource.TypeSubagent, scope) {
				plan.Tasks = append(plan.Tasks, InstallTask{
					Provider: prov,
					Resource: s,
					Scope:    scope,
				})
			}
		}
	}

	return plan
}

// BuildPlanFromProfile resolves a profile into its concrete resources and
// creates an install plan.
func BuildPlanFromProfile(
	profile resource.Profile,
	reg *resource.Registry,
	providers []provider.Provider,
	scope provider.Scope,
) *InstallPlan {
	skills, rules, mcps, subagents := resource.ResolveProfile(profile, reg)
	return BuildPlan(providers, skills, rules, mcps, subagents, scope)
}

// Execute runs all tasks in the plan and returns results.
func (p *InstallPlan) Execute(projectRoot, homeDir string, envOverrides map[string]string) []adapter.InstallResult {
	results := make([]adapter.InstallResult, 0, len(p.Tasks))

	for _, task := range p.Tasks {
		a := adapter.ForProvider(task.Provider, envOverrides)
		var res adapter.InstallResult

		switch r := task.Resource.(type) {
		case resource.Skill:
			res = a.InstallSkill(r, task.Scope, projectRoot, homeDir)
		case resource.Rule:
			res = a.InstallRule(r, task.Scope, projectRoot, homeDir)
		case resource.MCPServer:
			res = a.InstallMCP(r, task.Scope, projectRoot, homeDir)
		case resource.Subagent:
			res = a.InstallSubagent(r, task.Scope, projectRoot, homeDir)
		default:
			res = adapter.InstallResult{
				Provider: task.Provider,
				Resource: task.Resource,
				Scope:    task.Scope,
				Skipped:  true,
				Note:     "unknown resource type",
			}
		}

		results = append(results, res)
	}

	return results
}
