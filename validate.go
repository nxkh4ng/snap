package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

func canType(typeNames []string) func(string) error {
	return func(input string) error {
		input = strings.ToLower(strings.TrimSpace(input))
		if slices.Contains(typeNames, input) {
			return nil
		}
		return fmt.Errorf("invalid commit type: %s\n(allowed: %s)", input, strings.Join(typeNames, ", "))
	}
}

func canScope(scopes []string, required bool) func(string) error {
	return func(input string) error {
		if strings.TrimSpace(input) == "" {
			if required {
				return fmt.Errorf("scope is required")
			}
			return nil
		}

		if len(scopes) > 0 && !slices.Contains(scopes, input) {
			return fmt.Errorf("invalid scope: %s\n(allowed: %s)",
				input, strings.Join(scopes, ", "))
		}

		return nil
	}
}

func canSubject(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("subject is required")
	}
	if unicode.IsUpper(rune(input[0])) {
		return fmt.Errorf("subject must not start with uppercase")
	}
	if strings.HasSuffix(input, ".") {
		return fmt.Errorf("subject must not end with a period")
	}
	return nil
}
