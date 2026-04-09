package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode"

	"charm.land/huh/v2"
)

var (
	commitType     string
	scope          string
	subject        string
	desc           string
	footer         string
	breakingChange string
	ticketID       string
	ticketKeyWord  string
	confirmed      bool

	ticketKeyWords = []string{
		"Closes",
		"Fixes",
		"Refs",
	}

	types = []string{
		"feat",
		"fix",
		"chore",
		"docs",
		"style",
		"refactor",
		"perf",
		"test",
		"build",
		"ci",
	}

	typeDescriptions = map[string]string{
		"feat":     "A new feature",
		"fix":      "A bug fix",
		"chore":    "Build process or auxiliary tool changes",
		"docs":     "Documentation only changes",
		"style":    "Markup, white-space, formatting, missing semi-colons...",
		"refactor": "A code change that neither fixes a bug nor adds a feature",
		"perf":     "A code change that improves performance",
		"test":     "Adding missing tests",
		"build":    "Changes that affect the build system or external dependencies",
		"ci":       "CI related changes",
	}
)

func isCommitType(input string) error {
	input = strings.ToLower(strings.TrimSpace(input))
	if slices.Contains(types, input) {
		return nil
	}
	return fmt.Errorf("invalid commit type: %s\n(allowed: %s)", input, strings.Join(types, ", "))
}

func isSubject(input string) error {
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

func buildCommitMsg() string {
	sc := scope
	if sc != "" {
		sc = fmt.Sprintf("(%s)", sc)
	}

	bang := ""
	if breakingChange != "" {
		bang = "!"
	}

	parts := []string{
		fmt.Sprintf("%s%s%s: %s", commitType, sc, bang, subject),
	}

	if desc != "" {
		parts = append(parts, desc)
	}

	var footerLines []string

	if breakingChange != "" {
		footerLines = append(footerLines, "BREAKING CHANGE: "+breakingChange)
	}

	if footer != "" {
		footerLines = append(footerLines, footer)
	}

	if ticketID != "" && ticketKeyWord != "" {
		footerLines = append(footerLines, ticketKeyWord+" "+ticketID)
	}

	if len(footerLines) > 0 {
		parts = append(parts, strings.Join(footerLines, "\n"))
	}

	return strings.Join(parts, "\n\n")
}

func main() {
	withDesc := flag.Bool("d", false, "add description")
	withFooter := flag.Bool("f", false, "add footer")
	withBreakingChange := flag.Bool("b", false, "add breaking change")
	withTicket := flag.Bool("t", false, "add ticket")
	flag.Parse()

	inputGroup := huh.NewGroup(
		huh.NewInput().
			Title("Type?").
			DescriptionFunc(func() string {
				if typeDesc, ok := typeDescriptions[commitType]; ok {
					return typeDesc
				}
				return "..."
			}, &commitType).
			Suggestions(types).
			Validate(isCommitType).
			Value(&commitType),

		huh.NewInput().
			Title("Scope?").
			Placeholder("optional").
			Value(&scope),

		huh.NewInput().
			Title("Subject?").
			CharLimit(100).
			Validate(isSubject).
			Value(&subject),
	)

	descGroup := huh.NewGroup(
		huh.NewText().
			Title("Description?").
			Placeholder("longer explanation (optional)").
			Value(&desc),
	).WithHideFunc(func() bool {
		return !*withDesc
	})

	breakingChangeGroup := huh.NewGroup(
		huh.NewInput().
			Title("Breaking Change?").
			Prompt("BREAKING CHANGE: ").
			Value(&breakingChange),
	).WithHideFunc(func() bool {
		return !*withBreakingChange
	})

	footerGroup := huh.NewGroup(
		huh.NewText().
			Title("Footer?").
			Placeholder("longer footer (optional)").
			Value(&footer),
	).WithHideFunc(func() bool {
		return !*withFooter
	})

	ticketGroup := huh.NewGroup(
		huh.NewInput().
			Title("Ticket ID?").
			Placeholder("e.g. COMP-123, #42").
			Value(&ticketID),

		huh.NewSelect[string]().
			Title("Ticket keyword?").
			Options(huh.NewOptions(ticketKeyWords...)...).
			Value(&ticketKeyWord),
	).WithHideFunc(func() bool {
		return !*withTicket
	})

	form := huh.NewForm(
		inputGroup,
		descGroup,
		breakingChangeGroup,
		footerGroup,
		ticketGroup,
	)

	if err := form.Run(); err != nil {
		fmt.Println("Aborted.")
		return
	}

	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Commit this message?").
				DescriptionFunc(func() string {
					return buildCommitMsg()
				}, &subject).
				Value(&confirmed),
		),
	)

	if err := confirmForm.Run(); err != nil {
		fmt.Println("Aborted.")
		return
	}

	if !confirmed {
		fmt.Println("Aborted.")
		return
	}

	fmt.Println(buildCommitMsg())

	cmd := exec.Command("git", "commit", "-m", buildCommitMsg())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", out)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
