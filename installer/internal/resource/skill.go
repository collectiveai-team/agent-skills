package resource

import (
	"os"
	"path/filepath"
)

type skillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Skill struct {
	Name string
	Path string
	Desc string
	Body []byte
}

func (s Skill) ResourceName() string        { return s.Name }
func (s Skill) ResourceType() Type          { return TypeSkill }
func (s Skill) ResourcePath() string        { return s.Path }
func (s Skill) ResourceDescription() string { return s.Desc }

func DiscoverSkills(root string) ([]Skill, error) {
	skillsDir := filepath.Join(root, "registry", "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(skillsDir, entry.Name())
		markerPath := filepath.Join(candidate, "SKILL.md")
		data, err := os.ReadFile(markerPath)
		if err != nil {
			continue
		}
		var fm skillFrontmatter
		_, parseErr := ParseFrontmatter(data, &fm)
		desc := fm.Description
		if parseErr != nil {
			desc = ""
		}
		skills = append(skills, Skill{
			Name: entry.Name(),
			Path: candidate,
			Desc: desc,
		})
	}

	return skills, nil
}
