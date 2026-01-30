package skills

import (
	"os"
	"path/filepath"
)

type Skill struct {
	Name string
	Path string
}

func DiscoverSkills(root string) ([]Skill, error) {
	skillsDir := filepath.Join(root, "skills")
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
		if _, err := os.Stat(filepath.Join(candidate, "SKILL.md")); err != nil {
			continue
		}
		skills = append(skills, Skill{Name: entry.Name(), Path: candidate})
	}

	return skills, nil
}
