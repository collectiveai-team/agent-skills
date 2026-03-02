package resource

import (
	"os"
	"path/filepath"
)

type subagentFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Subagent struct {
	DirName string
	Path    string
	Desc    string
	Version string
	Body    []byte
}

func (s Subagent) ResourceName() string        { return s.DirName }
func (s Subagent) ResourceType() Type          { return TypeSubagent }
func (s Subagent) ResourcePath() string        { return s.Path }
func (s Subagent) ResourceDescription() string { return s.Desc }

func DiscoverSubagents(root string) ([]Subagent, error) {
	subDir := filepath.Join(root, "registry", "subagents")
	entries, err := os.ReadDir(subDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var subagents []Subagent
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(subDir, entry.Name())
		markerPath := filepath.Join(candidate, "SUBAGENT.md")
		data, err := os.ReadFile(markerPath)
		if err != nil {
			continue
		}
		var fm subagentFrontmatter
		body, parseErr := ParseFrontmatter(data, &fm)
		if parseErr != nil {
			continue
		}
		subagents = append(subagents, Subagent{
			DirName: entry.Name(),
			Path:    candidate,
			Desc:    fm.Description,
			Version: fm.Version,
			Body:    body,
		})
	}

	return subagents, nil
}
