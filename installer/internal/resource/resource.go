package resource

// Type represents the kind of resource in the registry.
type Type int

const (
	TypeSkill Type = iota
	TypeRule
	TypeMCPServer
	TypeSubagent
	TypeProfile
)

func (t Type) String() string {
	switch t {
	case TypeSkill:
		return "skill"
	case TypeRule:
		return "rule"
	case TypeMCPServer:
		return "mcp-server"
	case TypeSubagent:
		return "subagent"
	case TypeProfile:
		return "profile"
	default:
		return "unknown"
	}
}

// Resource is the common interface implemented by all registry entries.
type Resource interface {
	ResourceName() string
	ResourceType() Type
	ResourcePath() string
	ResourceDescription() string
}
