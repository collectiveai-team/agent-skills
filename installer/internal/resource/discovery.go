package resource

// Registry holds all discovered resources from the registry directory tree.
type Registry struct {
	Skills    []Skill
	Rules     []Rule
	MCPs      []MCPServer
	Subagents []Subagent
	Profiles  []Profile
}

// Discover scans the registry root for all resource types and returns a
// populated Registry. Individual discover errors for optional types (rules,
// MCP servers, subagents, profiles) are ignored — only skills are required.
func Discover(root string) (*Registry, error) {
	skills, err := DiscoverSkills(root)
	if err != nil {
		return nil, err
	}

	rules, _ := DiscoverRules(root)
	mcps, _ := DiscoverMCPServers(root)
	subagents, _ := DiscoverSubagents(root)
	profiles, _ := DiscoverProfiles(root)

	return &Registry{
		Skills:    skills,
		Rules:     rules,
		MCPs:      mcps,
		Subagents: subagents,
		Profiles:  profiles,
	}, nil
}

// ResolveProfile expands a Profile into concrete resources by matching names
// against the registry.
func ResolveProfile(profile Profile, reg *Registry) (skills []Skill, rules []Rule, mcps []MCPServer, subagents []Subagent) {
	skillIdx := make(map[string]Skill, len(reg.Skills))
	for _, s := range reg.Skills {
		skillIdx[s.Name] = s
	}
	ruleIdx := make(map[string]Rule, len(reg.Rules))
	for _, r := range reg.Rules {
		ruleIdx[r.Name] = r
	}
	mcpIdx := make(map[string]MCPServer, len(reg.MCPs))
	for _, m := range reg.MCPs {
		mcpIdx[m.DirName] = m
	}
	subIdx := make(map[string]Subagent, len(reg.Subagents))
	for _, s := range reg.Subagents {
		subIdx[s.DirName] = s
	}

	for _, name := range profile.Skills {
		if s, ok := skillIdx[name]; ok {
			skills = append(skills, s)
		}
	}
	for _, name := range profile.Rules {
		if r, ok := ruleIdx[name]; ok {
			rules = append(rules, r)
		}
	}
	for _, name := range profile.MCPServers {
		if m, ok := mcpIdx[name]; ok {
			mcps = append(mcps, m)
		}
	}
	for _, name := range profile.Subagents {
		if s, ok := subIdx[name]; ok {
			subagents = append(subagents, s)
		}
	}
	return
}
