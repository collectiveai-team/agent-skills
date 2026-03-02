package resource

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type MCPServerEntry struct {
	ID      string            `yaml:"id"`
	Name    string            `yaml:"name"`
	Command string            `yaml:"command"`
	Args    []string          `yaml:"args"`
	Env     map[string]string `yaml:"env,omitempty"`
}

type mcpManifest struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Version     string           `yaml:"version"`
	Servers     []MCPServerEntry `yaml:"servers"`
}

type MCPServer struct {
	DirName string
	Path    string
	Desc    string
	Version string
	Servers []MCPServerEntry
}

func (m MCPServer) ResourceName() string        { return m.DirName }
func (m MCPServer) ResourceType() Type          { return TypeMCPServer }
func (m MCPServer) ResourcePath() string        { return m.Path }
func (m MCPServer) ResourceDescription() string { return m.Desc }

func DiscoverMCPServers(root string) ([]MCPServer, error) {
	mcpDir := filepath.Join(root, "registry", "mcp-servers")
	entries, err := os.ReadDir(mcpDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var servers []MCPServer
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(mcpDir, entry.Name())
		markerPath := filepath.Join(candidate, "MCP.yaml")
		data, err := os.ReadFile(markerPath)
		if err != nil {
			continue
		}
		var m mcpManifest
		if err := yaml.Unmarshal(data, &m); err != nil {
			continue
		}
		servers = append(servers, MCPServer{
			DirName: entry.Name(),
			Path:    candidate,
			Desc:    m.Description,
			Version: m.Version,
			Servers: m.Servers,
		})
	}

	return servers, nil
}
