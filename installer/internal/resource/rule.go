package resource

import (
	"os"
	"path/filepath"
)

type ruleFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Rule struct {
	Name    string
	Path    string
	Desc    string
	Version string
	Body    []byte
}

func (r Rule) ResourceName() string        { return r.Name }
func (r Rule) ResourceType() Type          { return TypeRule }
func (r Rule) ResourcePath() string        { return r.Path }
func (r Rule) ResourceDescription() string { return r.Desc }

func DiscoverRules(root string) ([]Rule, error) {
	rulesDir := filepath.Join(root, "registry", "rules")
	entries, err := os.ReadDir(rulesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var rules []Rule
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(rulesDir, entry.Name())
		markerPath := filepath.Join(candidate, "RULE.md")
		data, err := os.ReadFile(markerPath)
		if err != nil {
			continue
		}
		var fm ruleFrontmatter
		body, parseErr := ParseFrontmatter(data, &fm)
		if parseErr != nil {
			continue
		}
		rules = append(rules, Rule{
			Name:    entry.Name(),
			Path:    candidate,
			Desc:    fm.Description,
			Version: fm.Version,
			Body:    body,
		})
	}

	return rules, nil
}
