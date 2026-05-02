package internal

import "strings"

func FormatCommitMsg(commitType, scope, summary, description, breakingChange string) string {
	var b strings.Builder

	commitType = strings.TrimSpace(commitType)
	scope = strings.ToLower(strings.TrimSpace(scope))
	summary = strings.TrimSpace(summary)
	description = strings.TrimSpace(description)
	breakingChange = strings.TrimSpace(breakingChange)

	hasBreakingChange := strings.Contains(commitType, "!")
	if hasBreakingChange {
		commitType = strings.TrimSuffix(commitType, "!")
	}

	// Header
	b.WriteString(commitType)
	if scope != "" {
		b.WriteString("(")
		b.WriteString(scope)
		b.WriteString(")")
	}
	if hasBreakingChange {
		b.WriteString("!")
	}
	b.WriteString(": ")
	b.WriteString(summary)

	// Description
	if description != "" {
		b.WriteString("\n\n")
		b.WriteString(description)
	}

	// BREAKING CHANGES
	if breakingChange != "" {
		b.WriteString("\n\nBREAKING CHANGE: ")
		b.WriteString(breakingChange)
	}

	return b.String()
}

func ParseCommitMsg(msg string) (commitType, scope, summary, description, breakingChange string, hasBreaking bool) {
	parts := strings.SplitN(msg, "\n\n", 2)
	header := strings.TrimSpace(parts[0])
	body := ""
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}

	headerParts := strings.SplitN(header, ": ", 2)
	if len(headerParts) == 2 {
		typeScope := headerParts[0]
		summary = headerParts[1]

		hasBreaking = strings.Contains(typeScope, "!")
		typeScope = strings.TrimSuffix(typeScope, "!")

		if before, after, ok := strings.Cut(typeScope, "("); ok {
			commitType = before
			scope = strings.TrimSuffix(after, ")")
		} else {
			commitType = typeScope
		}

		if hasBreaking {
			commitType = commitType + "!"
		}
	}

	if before, after, ok := strings.Cut(body, "BREAKING CHANGE:"); ok {
		description = strings.TrimSpace(before)
		breakingChange = strings.TrimSpace(after)
	} else {
		description = body
	}

	return commitType, scope, summary, description, breakingChange, hasBreaking
}
