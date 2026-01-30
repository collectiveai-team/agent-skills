package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"skill-installer/internal/agents"
	"skill-installer/internal/fs"
	"skill-installer/internal/skills"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateSelectAgents state = iota
	stateSelectScope
	stateConfirm
	stateCopy
	statePromptOverwrite
	stateDone
	stateError
)

type copyTask struct {
	Agent agents.Agent
	Dest  string
	Skill skills.Skill
}

type copyResult struct {
	Task copyTask
	Err  error
	Note string
}

type copyStepMsg struct {
	Err      error
	Conflict bool
	Task     copyTask
}

type Model struct {
	state state

	agents        []agents.Agent
	selected      []bool
	cursor        int
	selectedScope agents.Scope

	projectRoot string
	homeDir     string
	tempDir     string

	skills []skills.Skill

	copyTasks   []copyTask
	copyIndex   int
	results     []copyResult
	conflict    *copyTask
	conflictMsg string

	warning string
	errMsg  string
}

func NewModel() Model {
	projectRoot, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	agentsList := agents.AllAgents()
	selected := make([]bool, len(agentsList))

	return Model{
		state:         stateSelectAgents,
		agents:        agentsList,
		selected:      selected,
		cursor:        0,
		selectedScope: agents.ScopeProject,
		projectRoot:   projectRoot,
		homeDir:       homeDir,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadSkillsCmd()
}

func (m Model) loadSkillsCmd() tea.Cmd {
	return func() tea.Msg {
		tempDir, err := os.MkdirTemp("", "skill-installer-*")
		if err != nil {
			return errMsg{err: err}
		}

		if err := cloneRepo(repoURL(), tempDir); err != nil {
			_ = os.RemoveAll(tempDir)
			return errMsg{err: err}
		}

		skillsList, err := skills.DiscoverSkills(tempDir)
		if err != nil {
			_ = os.RemoveAll(tempDir)
			return errMsg{err: err}
		}
		if len(skillsList) == 0 {
			_ = os.RemoveAll(tempDir)
			return errMsg{err: fmt.Errorf("no skills found under %s", filepath.Join(tempDir, "skills"))}
		}
		return skillsMsg{skills: skillsList, tempDir: tempDir}
	}
}

type skillsMsg struct {
	skills  []skills.Skill
	tempDir string
}

type errMsg struct {
	err error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case skillsMsg:
		m.skills = msg.skills
		m.tempDir = msg.tempDir
		return m, nil
	case errMsg:
		m.state = stateError
		m.errMsg = msg.err.Error()
		m.cleanupTempDir()
		return m, nil
	case copyStepMsg:
		if msg.Conflict {
			m.state = statePromptOverwrite
			m.conflict = &msg.Task
			m.conflictMsg = "Destination already exists"
			return m, nil
		}
		m.results = append(m.results, copyResult{Task: msg.Task, Err: msg.Err})
		m.copyIndex++
		if m.copyIndex >= len(m.copyTasks) {
			m.state = stateDone
			m.cleanupTempDir()
			return m, nil
		}
		return m, m.nextCopyCmd()
	}

	switch m.state {
	case stateSelectAgents:
		return m.updateSelectAgents(msg)
	case stateSelectScope:
		return m.updateSelectScope(msg)
	case stateConfirm:
		return m.updateConfirm(msg)
	case stateCopy:
		return m, nil
	case statePromptOverwrite:
		return m.updatePromptOverwrite(msg)
	case stateDone:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "q" || key.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}
		return m, nil
	case stateError:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "q" || key.String() == "ctrl+c" {
				m.cleanupTempDir()
				return m, tea.Quit
			}
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) updateSelectAgents(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.agents)-1 {
			m.cursor++
		}
	case " ":
		m.selected[m.cursor] = !m.selected[m.cursor]
	case "a":
		allSelected := true
		for _, s := range m.selected {
			if !s {
				allSelected = false
				break
			}
		}
		for i := range m.selected {
			m.selected[i] = !allSelected
		}
	case "enter":
		if !m.hasSelections() {
			m.warning = "Select at least one agent"
			return m, nil
		}
		m.warning = ""
		m.state = stateSelectScope
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) updateSelectScope(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "left", "h", "right", "l", "tab":
		if m.selectedScope == agents.ScopeProject {
			m.selectedScope = agents.ScopeGlobal
		} else {
			m.selectedScope = agents.ScopeProject
		}
	case "enter":
		if m.selectedScope == agents.ScopeProject && !fs.IsGitRepo(m.projectRoot) {
			m.warning = "Project scope selected, but no .git directory found. Falling back to global."
			m.selectedScope = agents.ScopeGlobal
		}
		m.state = stateConfirm
		return m, nil
	case "esc":
		m.state = stateSelectAgents
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "enter":
		m.buildCopyTasks()
		if len(m.copyTasks) == 0 {
			m.state = stateError
			m.errMsg = "No copy tasks created."
			return m, nil
		}
		m.state = stateCopy
		m.copyIndex = 0
		m.results = nil
		return m, m.nextCopyCmd()
	case "esc":
		m.state = stateSelectScope
		return m, nil
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) updatePromptOverwrite(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.conflict == nil {
		m.state = stateCopy
		return m, m.nextCopyCmd()
	}

	switch key.String() {
	case "y":
		task := *m.conflict
		m.conflict = nil
		m.state = stateCopy
		return m, m.overwriteCmd(task)
	case "n":
		task := *m.conflict
		m.conflict = nil
		m.state = stateCopy
		m.results = append(m.results, copyResult{Task: task, Note: "skipped"})
		m.copyIndex++
		if m.copyIndex >= len(m.copyTasks) {
			m.state = stateDone
			m.cleanupTempDir()
			return m, nil
		}
		return m, m.nextCopyCmd()
	case "q", "ctrl+c":
		m.cleanupTempDir()
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) buildCopyTasks() {
	m.copyTasks = nil
	for i, agent := range m.agents {
		if !m.selected[i] {
			continue
		}
		destRoot, ok := agents.ResolveDestination(agent, m.selectedScope, m.projectRoot, m.homeDir)
		if !ok {
			continue
		}
		for _, skill := range m.skills {
			dest := filepath.Join(destRoot, skill.Name)
			m.copyTasks = append(m.copyTasks, copyTask{Agent: agent, Dest: dest, Skill: skill})
		}
	}
}

func (m Model) nextCopyCmd() tea.Cmd {
	if m.copyIndex >= len(m.copyTasks) {
		return nil
	}
	current := m.copyTasks[m.copyIndex]
	return func() tea.Msg {
		exists, err := fs.PathExists(current.Dest)
		if err != nil {
			return copyStepMsg{Err: err, Task: current}
		}
		if exists {
			return copyStepMsg{Conflict: true, Task: current}
		}
		if err := fs.EnsureDir(filepath.Dir(current.Dest)); err != nil {
			return copyStepMsg{Err: err, Task: current}
		}
		err = skills.CopyDir(current.Skill.Path, current.Dest)
		return copyStepMsg{Err: err, Task: current}
	}
}

func (m Model) overwriteCmd(task copyTask) tea.Cmd {
	return func() tea.Msg {
		if err := skills.RemoveIfExists(task.Dest); err != nil {
			return copyStepMsg{Err: err, Task: task}
		}
		if err := fs.EnsureDir(filepath.Dir(task.Dest)); err != nil {
			return copyStepMsg{Err: err, Task: task}
		}
		err := skills.CopyDir(task.Skill.Path, task.Dest)
		return copyStepMsg{Err: err, Task: task}
	}
}

func (m Model) hasSelections() bool {
	for _, s := range m.selected {
		if s {
			return true
		}
	}
	return false
}

func (m Model) View() string {
	switch m.state {
	case stateSelectAgents:
		return m.viewSelectAgents()
	case stateSelectScope:
		return m.viewSelectScope()
	case stateConfirm:
		return m.viewConfirm()
	case stateCopy:
		return m.viewCopy()
	case statePromptOverwrite:
		return m.viewPromptOverwrite()
	case stateDone:
		return m.viewDone()
	case stateError:
		return m.viewError()
	default:
		return ""
	}
}

func (m Model) viewSelectAgents() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Select coding agents"))
	b.WriteString("\n\n")

	for i, agent := range m.agents {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if m.selected[i] {
			checked = "x"
		}
		line := fmt.Sprintf("%s [%s] %s", cursor, checked, agent.Name)
		if m.cursor == i {
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
	b.WriteString(helpStyle().Render("Space: toggle  A: toggle all  Enter: next  Q: quit"))
	return b.String()
}

func (m Model) viewSelectScope() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Select scope"))
	b.WriteString("\n\n")

	projectLabel := "Project"
	globalLabel := "Global"
	if m.selectedScope == agents.ScopeProject {
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

func (m Model) viewConfirm() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Confirm copy"))
	b.WriteString("\n\n")

	selectedAgents := m.selectedAgentNames()
	if len(selectedAgents) == 0 {
		selectedAgents = []string{"(none)"}
	}

	b.WriteString("Agents:\n")
	for _, name := range selectedAgents {
		b.WriteString("- ")
		b.WriteString(name)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Scope: %s\n", m.scopeLabel()))
	b.WriteString("\n")
	b.WriteString("Skills:\n")
	for _, skill := range m.skills {
		b.WriteString("- ")
		b.WriteString(skill.Name)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Enter: start copy  Esc: back  Q: quit"))
	return b.String()
}

func (m Model) viewCopy() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Copying skills"))
	b.WriteString("\n\n")

	if m.copyIndex < len(m.copyTasks) {
		current := m.copyTasks[m.copyIndex]
		b.WriteString(fmt.Sprintf("%s -> %s\n", current.Skill.Name, current.Agent.Name))
		b.WriteString(fmt.Sprintf("Destination: %s\n", current.Dest))
	} else {
		b.WriteString("Finalizing...\n")
	}

	b.WriteString("\n")
	current := m.copyIndex + 1
	if current > len(m.copyTasks) {
		current = len(m.copyTasks)
	}
	b.WriteString(progressStyle().Render(fmt.Sprintf("%d/%d", current, len(m.copyTasks))))
	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Q: quit"))
	return b.String()
}

func (m Model) viewPromptOverwrite() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Overwrite existing skill?"))
	b.WriteString("\n\n")
	if m.conflict != nil {
		b.WriteString(fmt.Sprintf("Skill: %s\n", m.conflict.Skill.Name))
		b.WriteString(fmt.Sprintf("Agent: %s\n", m.conflict.Agent.Name))
		b.WriteString(fmt.Sprintf("Destination: %s\n", m.conflict.Dest))
	}
	if m.conflictMsg != "" {
		b.WriteString("\n")
		b.WriteString(warnStyle().Render(m.conflictMsg))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(helpStyle().Render("Y: overwrite  N: skip  Q: quit"))
	return b.String()
}

func (m Model) viewDone() string {
	var b strings.Builder
	b.WriteString(titleStyle().Render("Done"))
	b.WriteString("\n\n")

	if len(m.results) == 0 {
		b.WriteString("No results to show.\n")
	} else {
		for _, res := range m.results {
			status := "ok"
			if res.Note == "skipped" {
				status = "skipped"
			} else if res.Err != nil {
				status = "error"
			}
			line := fmt.Sprintf("[%s] %s -> %s", status, res.SkillName(), res.AgentName())
			b.WriteString(line)
			if res.Err != nil {
				b.WriteString("\n")
				b.WriteString(indentStyle().Render(res.Err.Error()))
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

func (m *Model) cleanupTempDir() {
	if m.tempDir == "" {
		return
	}
	_ = os.RemoveAll(m.tempDir)
	m.tempDir = ""
}

func repoURL() string {
	return "https://github.com/collectiveai-team/agent-skills.git"
}

func cloneRepo(url string, dest string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func (m Model) selectedAgentNames() []string {
	var names []string
	for i, agent := range m.agents {
		if m.selected[i] {
			names = append(names, agent.Name)
		}
	}
	return names
}

func (m Model) scopeLabel() string {
	if m.selectedScope == agents.ScopeProject {
		return "Project"
	}
	return "Global"
}

func (r copyResult) AgentName() string {
	return r.Task.Agent.Name
}

func (r copyResult) SkillName() string {
	return r.Task.Skill.Name
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}

func highlightStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)
}

func warnStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
}

func helpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
}

func progressStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
}

func indentStyle() lipgloss.Style {
	return lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("196"))
}
