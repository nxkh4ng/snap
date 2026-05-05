package utils

type ValidationConfig struct {
	SummaryMaxLen       int
	ScopeMaxLen         int
	ScopeRequired       bool
	DescriptionRequired bool
}

var DefaultCommitTypes = map[string]string{
	"feat":     "A new feature",
	"fix":      "A bug fix",
	"docs":     "Documentation only changes",
	"style":    "Formatting, white-space, missing semi-colons,...",
	"refactor": "Code changes that neither fix bugs nor add features",
	"perf":     "Code changes that improve performance",
	"test":     "Adding missing tests or correcting existing tests",
	"build":    "Changes that affect the build system or external dependencies",
	"ci":       "Changes to our CI configuration files and scripts",
	"chore":    "Other changes that don't modify src or test files",
	"revert":   "Reverts a previous commit",
}
