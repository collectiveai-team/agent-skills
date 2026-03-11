package ui

import (
	"fmt"
	"strings"

	"installer/internal/resource"
)

func (m Model) viewLoading() string {
	return titleStyle().Render("Agent Setup") + "\n\n" +
		progressStyle().Render("Loading registry...") + "\n\n" +
		helpStyle().Render("Q: quit")
}

func (m Model) viewSelectMode() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Agent Setup"))
	b.WriteString("\n\n")

	modes := []struct {
		name string
		desc string
	}{
		{"Install profile", "Install a curated set of skills, rules, MCP servers, and subagents"},
		{"Install resources", "Pick individual resources to install"},
	}

	for i, mode := range modes {
		cursor := " "
		if m.modeCursor == i {
			cursor = ">"
		}
		line := fmt.Sprintf("%s %s", cursor, mode.name)
		if m.modeCursor == i {
			line = highlightStyle().Render(line)
			b.WriteString(line)
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(dimStyle().Render(mode.desc))
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	if m.warning != "" {
		b.WriteString("\n")
		b.WriteString(warnStyle().Render(m.warning))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Enter: select  Q: quit"))
	return b.String()
}

func (m Model) viewSelectProfile() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Select profile"))
	b.WriteString("\n\n")

	for i, profile := range m.registry.Profiles {
		cursor := " "
		if m.profileCursor == i {
			cursor = ">"
		}
		line := fmt.Sprintf("%s %s", cursor, profile.DirName)
		if m.profileCursor == i {
			line = highlightStyle().Render(line)
			b.WriteString(line)
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(dimStyle().Render(profile.Desc))
			b.WriteString("\n")
			// Show contents
			if len(profile.Skills) > 0 {
				b.WriteString(fmt.Sprintf("  Skills: %s\n", dimStyle().Render(strings.Join(profile.Skills, ", "))))
			}
			if len(profile.Rules) > 0 {
				b.WriteString(fmt.Sprintf("  Rules: %s\n", dimStyle().Render(strings.Join(profile.Rules, ", "))))
			}
			if len(profile.MCPServers) > 0 {
				b.WriteString(fmt.Sprintf("  MCP: %s\n", dimStyle().Render(strings.Join(profile.MCPServers, ", "))))
			}
			if len(profile.Subagents) > 0 {
				b.WriteString(fmt.Sprintf("  Subagents: %s\n", dimStyle().Render(strings.Join(profile.Subagents, ", "))))
			}
		} else {
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Enter: select  Esc: back  Q: quit"))
	return b.String()
}

func (m Model) viewSelectResources() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Select resources"))
	b.WriteString("\n\n")

	// Tab bar
	for i := resourceTab(0); i < tabCount; i++ {
		label := i.String()
		if i == m.activeTab {
			label = tabActiveStyle().Render(label)
		} else {
			label = tabInactiveStyle().Render(label)
		}
		b.WriteString(label)
		if i < tabCount-1 {
			b.WriteString("  ")
		}
	}
	b.WriteString("\n\n")

	// Items for current tab
	items := m.currentTabItems()
	sel := m.currentTabSelected()

	if len(items) == 0 {
		b.WriteString(dimStyle().Render("(none available)"))
		b.WriteString("\n")
	} else {
		for i, name := range items {
			cursor := " "
			if m.resourceCursor == i {
				cursor = ">"
			}
			checked := " "
			if sel[i] {
				checked = "x"
			}
			line := fmt.Sprintf("%s [%s] %s", cursor, checked, name)
			if m.resourceCursor == i {
				line = highlightStyle().Render(line)
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	if m.warning != "" {
		b.WriteString("\n")
		b.WriteString(warnStyle().Render(m.warning))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Tab: switch tab  Space: toggle  A: toggle all  Enter: next  Esc: back  Q: quit"))
	return b.String()
}

func (m Model) viewSelectProviders() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Select providers"))
	b.WriteString("\n\n")

	for i, prov := range m.providers {
		cursor := " "
		if m.providerCursor == i {
			cursor = ">"
		}
		checked := " "
		if m.providerSelected[i] {
			checked = "x"
		}
		line := fmt.Sprintf("%s [%s] %s", cursor, checked, prov.Name)
		if m.providerCursor == i {
			line = highlightStyle().Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	if m.warning != "" {
		b.WriteString("\n")
		b.WriteString(warnStyle().Render(m.warning))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Space: toggle  A: toggle all  Enter: next  Esc: back  Q: quit"))
	return b.String()
}

func (m Model) viewSelectScope() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Select scope"))
	b.WriteString("\n\n")

	projectLabel := "Project"
	globalLabel := "Global"
	if m.selectedScope == 0 { // ScopeProject
		projectLabel = highlightStyle().Render(projectLabel)
	} else {
		globalLabel = highlightStyle().Render(globalLabel)
	}

	b.WriteString(fmt.Sprintf("%s  %s\n", projectLabel, globalLabel))
	if m.warning != "" {
		b.WriteString("\n")
		b.WriteString(warnStyle().Render(m.warning))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Tab/Arrows: toggle  Enter: next  Esc: back  Q: quit"))
	return b.String()
}

func (m Model) viewConfigureMCP() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Configure MCP server environment"))
	b.WriteString("\n\n")

	if len(m.mcpEnvKeys) == 0 {
		b.WriteString(dimStyle().Render("No environment variables to configure."))
		b.WriteString("\n")
	} else {
		for i, key := range m.mcpEnvKeys {
			cursor := " "
			if m.mcpEnvCursor == i {
				cursor = ">"
			}
			val := m.mcpEnvValues[key]
			line := fmt.Sprintf("%s %s = %s", cursor, key, val)
			if m.mcpEnvCursor == i {
				line = highlightStyle().Render(line)
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(dimStyle().Render("Values resolved from your environment. Press Enter to continue."))
	b.WriteString("\n\n")
	b.WriteString(helpStyle().Render("Enter: continue  Esc: back  Q: quit"))
	return b.String()
}

func (m Model) viewConfirm() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Confirm installation"))
	b.WriteString("\n\n")

	// Providers
	selectedProviders := m.selectedProviderNames()
	if len(selectedProviders) == 0 {
		selectedProviders = []string{"(none)"}
	}
	b.WriteString("Providers:\n")
	for _, name := range selectedProviders {
		b.WriteString("  - ")
		b.WriteString(name)
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("\nScope: %s\n", m.scopeLabel()))

	// Group tasks by resource type
	if m.installPlan != nil {
		skillNames := m.selectedResourceNames(resource.TypeSkill)
		ruleNames := m.selectedResourceNames(resource.TypeRule)
		mcpNames := m.selectedResourceNames(resource.TypeMCPServer)
		subNames := m.selectedResourceNames(resource.TypeSubagent)

		if len(skillNames) > 0 {
			b.WriteString("\nSkills:\n")
			for _, name := range skillNames {
				b.WriteString("  - ")
				b.WriteString(name)
				b.WriteString("\n")
			}
		}
		if len(ruleNames) > 0 {
			b.WriteString("\nRules:\n")
			for _, name := range ruleNames {
				b.WriteString("  - ")
				b.WriteString(name)
				b.WriteString("\n")
			}
		}
		if len(mcpNames) > 0 {
			b.WriteString("\nMCP Servers:\n")
			for _, name := range mcpNames {
				b.WriteString("  - ")
				b.WriteString(name)
				b.WriteString("\n")
			}
		}
		if len(subNames) > 0 {
			b.WriteString("\nSubagents:\n")
			for _, name := range subNames {
				b.WriteString("  - ")
				b.WriteString(name)
				b.WriteString("\n")
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Enter: install  Esc: back  Q: quit"))
	return b.String()
}

func (m Model) selectedResourceNames(resType resource.Type) []string {
	seen := make(map[string]bool)
	var names []string
	if m.installPlan == nil {
		return nil
	}
	for _, task := range m.installPlan.Tasks {
		if task.Resource.ResourceType() == resType {
			name := task.Resource.ResourceName()
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

func (m Model) viewInstall() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Installing..."))
	b.WriteString("\n\n")

	total := 0
	if m.installPlan != nil {
		total = len(m.installPlan.Tasks)
	}
	b.WriteString(progressStyle().Render(fmt.Sprintf("Processing %d tasks...", total)))
	b.WriteString("\n\n")
	b.WriteString(helpStyle().Render("Please wait..."))
	return b.String()
}

func (m Model) viewDone() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Done"))
	b.WriteString("\n\n")

	if len(m.results) == 0 {
		b.WriteString("No results to show.\n")
	} else {
		var okCount, skipCount, errCount int
		for _, res := range m.results {
			if res.Err != nil {
				errCount++
			} else if res.Skipped {
				skipCount++
			} else {
				okCount++
			}
		}

		b.WriteString(fmt.Sprintf("Installed: %s  Skipped: %s  Errors: %s\n\n",
			successStyle().Render(fmt.Sprintf("%d", okCount)),
			dimStyle().Render(fmt.Sprintf("%d", skipCount)),
			warnStyle().Render(fmt.Sprintf("%d", errCount)),
		))

		for _, res := range m.results {
			status := successStyle().Render("ok")
			if res.Skipped {
				status = dimStyle().Render("skip")
			} else if res.Err != nil {
				status = warnStyle().Render("err")
			}

			resType := res.Resource.ResourceType().String()
			line := fmt.Sprintf("[%s] %s %s -> %s",
				status,
				dimStyle().Render(resType),
				res.Resource.ResourceName(),
				res.Provider.Name,
			)
			b.WriteString(line)
			if res.Err != nil {
				b.WriteString("\n")
				b.WriteString(indentStyle().Render(res.Err.Error()))
			}
			if res.Skipped && res.Note != "" {
				b.WriteString(" ")
				b.WriteString(dimStyle().Render("(" + res.Note + ")"))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Q: quit"))
	return b.String()
}

func (m Model) viewError() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Error"))
	b.WriteString("\n\n")
	b.WriteString(warnStyle().Render(m.errMsg))
	b.WriteString("\n\n")
	b.WriteString(helpStyle().Render("Q: quit"))
	return b.String()
}
