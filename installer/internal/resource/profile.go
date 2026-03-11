package resource

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type profileManifest struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
	Skills      []string `yaml:"skills"`
	Rules       []string `yaml:"rules"`
	MCPServers  []string `yaml:"mcp_servers"`
	Subagents   []string `yaml:"subagents"`
}

type Profile struct {
	DirName    string
	Path       string
	Desc       string
	Version    string
	Skills     []string
	Rules      []string
	MCPServers []string
	Subagents  []string
}

func (p Profile) ResourceName() string        { return p.DirName }
func (p Profile) ResourceType() Type          { return TypeProfile }
func (p Profile) ResourcePath() string        { return p.Path }
func (p Profile) ResourceDescription() string { return p.Desc }

func DiscoverProfiles(root string) ([]Profile, error) {
	profDir := filepath.Join(root, "registry", "profiles")
	entries, err := os.ReadDir(profDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var profiles []Profile
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(profDir, entry.Name())
		markerPath := filepath.Join(candidate, "PROFILE.yaml")
		data, err := os.ReadFile(markerPath)
		if err != nil {
			continue
		}
		var m profileManifest
		if err := yaml.Unmarshal(data, &m); err != nil {
			continue
		}
		profiles = append(profiles, Profile{
			DirName:    entry.Name(),
			Path:       candidate,
			Desc:       m.Description,
			Version:    m.Version,
			Skills:     m.Skills,
			Rules:      m.Rules,
			MCPServers: m.MCPServers,
			Subagents:  m.Subagents,
		})
	}

	return profiles, nil
}
