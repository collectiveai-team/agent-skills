package resource

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

var frontmatterSep = []byte("---")

// ParseFrontmatter splits a markdown file with YAML frontmatter into the
// decoded frontmatter struct and the remaining markdown body.
// The file must start with "---\n".
func ParseFrontmatter(data []byte, dest interface{}) (body []byte, err error) {
	data = bytes.TrimLeft(data, "\n\r")
	if !bytes.HasPrefix(data, frontmatterSep) {
		return nil, fmt.Errorf("missing frontmatter separator")
	}

	rest := data[len(frontmatterSep):]
	idx := bytes.Index(rest, frontmatterSep)
	if idx < 0 {
		return nil, fmt.Errorf("missing closing frontmatter separator")
	}

	fmData := rest[:idx]
	body = bytes.TrimLeft(rest[idx+len(frontmatterSep):], "\n\r")

	if err := yaml.Unmarshal(fmData, dest); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}
	return body, nil
}
