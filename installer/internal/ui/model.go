package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"installer/internal/adapter"
	"installer/internal/fs"
	"installer/internal/installer"
	"installer/internal/provider"
	"installer/internal/resource"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateLoading state = iota
	stateSelectMode
	stateSelectProfile
	stateSelectResources
	stateSelectProviders
	stateSelectScope
	stateConfigureMCP
	stateConfirm
	stateInstall
	stateDone
	stateError
)

// installMode captures the top-level workflow path.
type installMode int

const (
	modeProfile   installMode = iota
	modeResources
)

// resourceTab enumerates the tabs in stateSelectResources.
type resourceTab int

const (
	tabSkills resourceTab = iota
	tabRules
	tabMCP
	tabSubagents
	tabCount // sentinel
)

func (t resourceTab) String() string {
	switch t {
	case tabSkills:
		return "Skills"
	case tabRules:
		return "Rules"
	case tabMCP:
		return "MCP Servers"
	case tabSubagents:
		return "Subagents"
	default:
		return ""
	}
}

// --- Messages ---

type registryMsg struct {
	registry *resource.Registry
	tempDir  string
}

type errMsg struct {
	err error
}

type installDoneMsg struct {
	results []adapter.InstallResult
}

// --- Model ---

type Model struct {
	state state

	// Mode selection
	mode       installMode
	modeCursor int

	// Provider selection
	providers        []provider.Provider
	providerSelected []bool
	providerCursor   int

	// Scope
	selectedScope provider.Scope

	// Registry data
	registry *resource.Registry
	tempDir  string

	// Profile selection
	profileCursor int

	// Resource selection (multi-tab)
	activeTab      resourceTab
	skillSelected  []bool
	ruleSelected   []bool
	mcpSelected    []bool
	subSelected    []bool
	resourceCursor int

	// MCP env config
	mcpEnvKeys    []string            // ordered list of "serverID:VARNAME"
	mcpEnvValues  map[string]string   // "serverID:VARNAME" -> resolved value
	mcpEnvCursor  int

	// Install results
	results     []adapter.InstallResult
	installPlan *installer.InstallPlan

	// Environment
	projectRoot string
	homeDir     string

	// UI state
	warning string
	errMsg  string
}

func NewModel() Model {
	projectRoot, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	provs := provider.AllProviders()
	selected := make([]bool, len(provs))

	return Model{
		state:            stateLoading,
		providers:        provs,
		providerSelected: selected,
		providerCursor:   0,
		selectedScope:    provider.ScopeProject,
		projectRoot:      projectRoot,
		homeDir:          homeDir,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadRegistryCmd()
}

func (m Model) loadRegistryCmd() tea.Cmd {
	return func() tea.Msg {
		tempDir, err := os.MkdirTemp("", "agent-setup-*")
		if err != nil {
			return errMsg{err: err}
		}

		if err := cloneRepo(repoURL(), tempDir); err != nil {
			_ = os.RemoveAll(tempDir)
			return errMsg{err: err}
		}

		reg, err := resource.Discover(tempDir)
		if err != nil {
			_ = os.RemoveAll(tempDir)
			return errMsg{err: err}
		}
		if len(reg.Skills) == 0 && len(reg.Rules) == 0 && len(reg.MCPs) == 0 && len(reg.Subagents) == 0 {
			_ = os.RemoveAll(tempDir)
			return errMsg{err: fmt.Errorf("no resources found in %s", filepath.Join(tempDir, "registry"))}
		}
		return registryMsg{registry: reg, tempDir: tempDir}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case registryMsg:
		m.registry = msg.registry
		m.tempDir = msg.tempDir
		m.state = stateSelectMode
		return m, nil
	case errMsg:
		m.state = stateError
		m.errMsg = msg.err.Error()
		m.cleanupTempDir()
		return m, nil
	case installDoneMsg:
		m.results = msg.results
		m.state = stateDone
		m.cleanupTempDir()
		return m, nil
	}

	switch m.state {
	case stateLoading:
		return m, nil
	case stateSelectMode:
		return m.updateSelectMode(msg)
	case stateSelectProfile:
		return m.updateSelectProfile(msg)
	case stateSelectResources:
		return m.updateSelectResources(msg)
	case stateSelectProviders:
		return m.updateSelectProviders(msg)
	case stateSelectScope:
		return m.updateSelectScope(msg)
	case stateConfigureMCP:
		return m.updateConfigureMCP(msg)
	case stateConfirm:
		return m.updateConfirm(msg)
	case stateInstall:
		return m, nil
	case stateDone, stateError:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == keyQ || key.String() == keyCtrlC {
				m.cleanupTempDir()
				return m, tea.Quit
			}
		}
		return m, nil
	default:
		return m, nil
	}
}

// --- State handlers ---

func (m Model) updateSelectMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case keyUp, keyK:
		if m.modeCursor > 0 {
			m.modeCursor--
		}
	case keyDown, keyJ:
		if m.modeCursor < 1 {
			m.modeCursor++
		}
	case keyEnter:
		switch m.modeCursor {
		case 0: // Install profile
			m.mode = modeProfile
			if len(m.registry.Profiles) == 0 {
				m.warning = "No profiles found in registry"
				return m, nil
			}
			m.warning = ""
			m.profileCursor = 0
			m.state = stateSelectProfile
		case 1: // Install resources
			m.mode = modeResources
			m.warning = ""
			m.initResourceSelection()
			m.state = stateSelectResources
		}
	case keyQ, keyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateSelectProfile(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	profiles := m.registry.Profiles
	switch key.String() {
	case keyUp, keyK:
		if m.profileCursor > 0 {
			m.profileCursor--
		}
	case keyDown, keyJ:
		if m.profileCursor < len(profiles)-1 {
			m.profileCursor++
		}
	case keyEnter:
		// Resolve profile into resources
		profile := profiles[m.profileCursor]
		skills, rules, mcps, subagents := resource.ResolveProfile(profile, m.registry)
		m.initResourceSelectionFrom(skills, rules, mcps, subagents)
		m.state = stateSelectProviders
	case keyEsc:
		m.state = stateSelectMode
	case keyQ, keyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) initResourceSelection() {
	m.skillSelected = make([]bool, len(m.registry.Skills))
	m.ruleSelected = make([]bool, len(m.registry.Rules))
	m.mcpSelected = make([]bool, len(m.registry.MCPs))
	m.subSelected = make([]bool, len(m.registry.Subagents))
	m.activeTab = tabSkills
	m.resourceCursor = 0
}

func (m *Model) initResourceSelectionFrom(skills []resource.Skill, rules []resource.Rule, mcps []resource.MCPServer, subagents []resource.Subagent) {
	// Rebuild selection arrays from the full registry, pre-selecting profile items
	m.skillSelected = make([]bool, len(m.registry.Skills))
	m.ruleSelected = make([]bool, len(m.registry.Rules))
	m.mcpSelected = make([]bool, len(m.registry.MCPs))
	m.subSelected = make([]bool, len(m.registry.Subagents))

	profileSkills := make(map[string]bool)
	for _, s := range skills {
		profileSkills[s.Name] = true
	}
	for i, s := range m.registry.Skills {
		m.skillSelected[i] = profileSkills[s.Name]
	}

	profileRules := make(map[string]bool)
	for _, r := range rules {
		profileRules[r.Name] = true
	}
	for i, r := range m.registry.Rules {
		m.ruleSelected[i] = profileRules[r.Name]
	}

	profileMCPs := make(map[string]bool)
	for _, m := range mcps {
		profileMCPs[m.DirName] = true
	}
	for i, mc := range m.registry.MCPs {
		m.mcpSelected[i] = profileMCPs[mc.DirName]
	}

	profileSubs := make(map[string]bool)
	for _, s := range subagents {
		profileSubs[s.DirName] = true
	}
	for i, s := range m.registry.Subagents {
		m.subSelected[i] = profileSubs[s.DirName]
	}
}

func (m Model) updateSelectResources(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	items := m.currentTabItems()

	switch key.String() {
	case keyTab:
		m.activeTab = (m.activeTab + 1) % tabCount
		m.resourceCursor = 0
	case keyUp, keyK:
		if m.resourceCursor > 0 {
			m.resourceCursor--
		}
	case keyDown, keyJ:
		if m.resourceCursor < len(items)-1 {
			m.resourceCursor++
		}
	case keySpace:
		if len(items) > 0 {
			m.toggleCurrentResource()
		}
	case keyA:
		m.toggleAllInCurrentTab()
	case keyEnter:
		if !m.hasAnyResourceSelected() {
			m.warning = "Select at least one resource"
			return m, nil
		}
		m.warning = ""
		m.state = stateSelectProviders
	case keyEsc:
		m.state = stateSelectMode
	case keyQ, keyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) currentTabItems() []string {
	switch m.activeTab {
	case tabSkills:
		names := make([]string, len(m.registry.Skills))
		for i, s := range m.registry.Skills {
			names[i] = s.Name
		}
		return names
	case tabRules:
		names := make([]string, len(m.registry.Rules))
		for i, r := range m.registry.Rules {
			names[i] = r.Name
		}
		return names
	case tabMCP:
		names := make([]string, len(m.registry.MCPs))
		for i, m := range m.registry.MCPs {
			names[i] = m.DirName
		}
		return names
	case tabSubagents:
		names := make([]string, len(m.registry.Subagents))
		for i, s := range m.registry.Subagents {
			names[i] = s.DirName
		}
		return names
	}
	return nil
}

func (m *Model) currentTabSelected() []bool {
	switch m.activeTab {
	case tabSkills:
		return m.skillSelected
	case tabRules:
		return m.ruleSelected
	case tabMCP:
		return m.mcpSelected
	case tabSubagents:
		return m.subSelected
	}
	return nil
}

func (m *Model) toggleCurrentResource() {
	sel := m.currentTabSelected()
	if m.resourceCursor < len(sel) {
		sel[m.resourceCursor] = !sel[m.resourceCursor]
	}
}

func (m *Model) toggleAllInCurrentTab() {
	sel := m.currentTabSelected()
	allSelected := true
	for _, s := range sel {
		if !s {
			allSelected = false
			break
		}
	}
	for i := range sel {
		sel[i] = !allSelected
	}
}

func (m Model) hasAnyResourceSelected() bool {
	for _, s := range m.skillSelected {
		if s {
			return true
		}
	}
	for _, s := range m.ruleSelected {
		if s {
			return true
		}
	}
	for _, s := range m.mcpSelected {
		if s {
			return true
		}
	}
	for _, s := range m.subSelected {
		if s {
			return true
		}
	}
	return false
}

func (m Model) updateSelectProviders(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case keyUp, keyK:
		if m.providerCursor > 0 {
			m.providerCursor--
		}
	case keyDown, keyJ:
		if m.providerCursor < len(m.providers)-1 {
			m.providerCursor++
		}
	case keySpace:
		m.providerSelected[m.providerCursor] = !m.providerSelected[m.providerCursor]
	case keyA:
		allSelected := true
		for _, s := range m.providerSelected {
			if !s {
				allSelected = false
				break
			}
		}
		for i := range m.providerSelected {
			m.providerSelected[i] = !allSelected
		}
	case keyEnter:
		if !m.hasProviderSelected() {
			m.warning = "Select at least one provider"
			return m, nil
		}
		m.warning = ""
		m.state = stateSelectScope
	case keyEsc:
		if m.mode == modeProfile {
			m.state = stateSelectProfile
		} else {
			m.state = stateSelectResources
		}
	case keyQ, keyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) hasProviderSelected() bool {
	for _, s := range m.providerSelected {
		if s {
			return true
		}
	}
	return false
}

func (m Model) updateSelectScope(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case keyLeft, keyH, keyRight, keyL, keyTab:
		if m.selectedScope == provider.ScopeProject {
			m.selectedScope = provider.ScopeGlobal
		} else {
			m.selectedScope = provider.ScopeProject
		}
	case keyEnter:
		if m.selectedScope == provider.ScopeProject && !fs.IsGitRepo(m.projectRoot) {
			m.warning = "No .git directory found. Falling back to global."
			m.selectedScope = provider.ScopeGlobal
		}
		// Check if any MCP servers are selected — go to configure screen
		if m.hasSelectedMCPServers() {
			m.buildMCPEnvConfig()
			m.state = stateConfigureMCP
		} else {
			m.buildInstallPlan()
			m.state = stateConfirm
		}
	case keyEsc:
		m.state = stateSelectProviders
	case keyQ, keyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) hasSelectedMCPServers() bool {
	for _, s := range m.mcpSelected {
		if s {
			return true
		}
	}
	return false
}

func (m *Model) buildMCPEnvConfig() {
	m.mcpEnvKeys = nil
	m.mcpEnvValues = make(map[string]string)

	for i, selected := range m.mcpSelected {
		if !selected {
			continue
		}
		mcp := m.registry.MCPs[i]
		for _, srv := range mcp.Servers {
			for k, v := range srv.Env {
				key := srv.ID + ":" + k
				m.mcpEnvKeys = append(m.mcpEnvKeys, key)
				// Resolve default value
				resolved := adapter.InterpolateEnv(v, nil)
				m.mcpEnvValues[key] = resolved
			}
		}
	}
	m.mcpEnvCursor = 0
}

func (m Model) updateConfigureMCP(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case keyUp, keyK:
		if m.mcpEnvCursor > 0 {
			m.mcpEnvCursor--
		}
	case keyDown, keyJ:
		if m.mcpEnvCursor < len(m.mcpEnvKeys)-1 {
			m.mcpEnvCursor++
		}
	case keyEnter:
		m.buildInstallPlan()
		m.state = stateConfirm
	case keyEsc:
		m.state = stateSelectScope
	case keyQ, keyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) buildInstallPlan() {
	var selectedProviders []provider.Provider
	for i, prov := range m.providers {
		if m.providerSelected[i] {
			selectedProviders = append(selectedProviders, prov)
		}
	}

	var skills []resource.Skill
	for i, s := range m.registry.Skills {
		if m.skillSelected[i] {
			skills = append(skills, s)
		}
	}

	var rules []resource.Rule
	for i, r := range m.registry.Rules {
		if m.ruleSelected[i] {
			rules = append(rules, r)
		}
	}

	var mcps []resource.MCPServer
	for i, mc := range m.registry.MCPs {
		if m.mcpSelected[i] {
			mcps = append(mcps, mc)
		}
	}

	var subs []resource.Subagent
	for i, s := range m.registry.Subagents {
		if m.subSelected[i] {
			subs = append(subs, s)
		}
	}

	m.installPlan = installer.BuildPlan(selectedProviders, skills, rules, mcps, subs, m.selectedScope)
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case keyEnter:
		if m.installPlan == nil || len(m.installPlan.Tasks) == 0 {
			m.state = stateError
			m.errMsg = "No install tasks to execute."
			return m, nil
		}
		m.state = stateInstall
		return m, m.executeInstallCmd()
	case keyEsc:
		if m.hasSelectedMCPServers() {
			m.state = stateConfigureMCP
		} else {
			m.state = stateSelectScope
		}
	case keyQ, keyCtrlC:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) executeInstallCmd() tea.Cmd {
	plan := m.installPlan
	projectRoot := m.projectRoot
	homeDir := m.homeDir
	envOverrides := m.mcpEnvValues

	return func() tea.Msg {
		results := plan.Execute(projectRoot, homeDir, envOverrides)
		return installDoneMsg{results: results}
	}
}

// --- View ---

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.viewLoading()
	case stateSelectMode:
		return m.viewSelectMode()
	case stateSelectProfile:
		return m.viewSelectProfile()
	case stateSelectResources:
		return m.viewSelectResources()
	case stateSelectProviders:
		return m.viewSelectProviders()
	case stateSelectScope:
		return m.viewSelectScope()
	case stateConfigureMCP:
		return m.viewConfigureMCP()
	case stateConfirm:
		return m.viewConfirm()
	case stateInstall:
		return m.viewInstall()
	case stateDone:
		return m.viewDone()
	case stateError:
		return m.viewError()
	default:
		return ""
	}
}

// --- Helpers ---

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

func (m Model) selectedProviderNames() []string {
	var names []string
	for i, prov := range m.providers {
		if m.providerSelected[i] {
			names = append(names, prov.Name)
		}
	}
	return names
}

func (m Model) scopeLabel() string {
	if m.selectedScope == provider.ScopeProject {
		return "Project"
	}
	return "Global"
}
