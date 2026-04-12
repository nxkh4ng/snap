package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
)

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

func scopesFromHistory() []string {
	cmd := exec.Command("git", "log", "--pretty=format:%s", "-50")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	seen := map[string]bool{}
	var scopes []string

	// parse scope from "type(scope): subject"
	re := regexp.MustCompile(`\(([^)]+)\)`)
	for line := range strings.SplitSeq(string(out), "\n") {
		if m := re.FindStringSubmatch(line); len(m) > 1 {
			s := m[1]
			if !seen[s] {
				seen[s] = true
				scopes = append(scopes, s)
			}
		}
	}
	return scopes
}

func runInitCmd() {
	isGlobal := false
	for _, arg := range os.Args[2:] {
		if arg == "--global" || arg == "-g" {
			isGlobal = true
			break
		}
	}

	if !isGlobal {
		fmt.Println("usage: snap init --global")
		return
	}

	path := globalConfigPath()

	if err := saveConfig(defaultConfig); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write config: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("config saved to %s\n", path)
}

func getTheme(name string) huh.ThemeFunc {
	themes := map[string]huh.ThemeFunc{
		"base":       huh.ThemeBase,
		"base16":     huh.ThemeBase16,
		"catppuccin": huh.ThemeCatppuccin,
		"dracula":    huh.ThemeDracula,
		"charm":      huh.ThemeCharm,
	}

	if theme, ok := themes[strings.ToLower(name)]; ok {
		return theme
	}

	return huh.ThemeBase16
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			runInitCmd()
			return
		}
	}
	initFlags()
	if handleFlags() {
		return
	}

	cfg := loadConfig()
	typeNames := make([]string, len(cfg.Types))
	typeDescs := make(map[string]string)
	for i, t := range cfg.Types {
		typeNames[i] = t.Name
		typeDescs[t.Name] = t.Description
	}
	ticketDescs := make(map[string]string)
	for _, kw := range cfg.TicketKeyWords {
		ticketDescs[kw.Name] = kw.Description
	}
	theme := getTheme(cfg.Theme)

	scopesHistory := scopesFromHistory()

	inputGroup := huh.NewGroup(
		huh.NewInput().
			Title("Type").
			DescriptionFunc(func() string {
				if typeDesc, ok := typeDescs[commitType]; ok {
					return typeDesc
				}
				return "..."
			}, &commitType).
			Suggestions(typeNames).
			Validate(canType(typeNames)).
			Value(&commitType),

		huh.NewInput().
			Title("Scope").
			Validate(canScope(cfg.Scopes, cfg.RequireScope)).
			Suggestions(scopesHistory).
			Value(&scope),

		huh.NewInput().
			Title("Subject").
			CharLimit(cfg.SubjectCharLimit).
			Validate(canSubject).
			Value(&subject),
	)

	descGroup := huh.NewGroup(
		huh.NewText().
			Title("Description").
			Placeholder("longer explanation (optional)").
			Value(&desc),
	).WithHideFunc(func() bool {
		return !*withDesc
	})

	breakingChangeGroup := huh.NewGroup(
		huh.NewInput().
			Title("BREAKING CHANGE").
			Value(&breakingChange),
	).WithHideFunc(func() bool {
		return !*withBreakingChange
	})

	footerGroup := huh.NewGroup(
		huh.NewText().
			Title("Footer").
			Placeholder("longer footer (optional)").
			Value(&footer),
	).WithHideFunc(func() bool {
		return !*withFooter
	})

	ticketGroup := huh.NewGroup(
		huh.NewInput().
			Title("Ticket ID").
			Placeholder("COMP-123, #42").
			Value(&ticketID),

		huh.NewSelect[string]().
			Title("Ticket keyword?").
			DescriptionFunc(func() string {
				if desc, ok := ticketDescs[ticketKeyWord]; ok {
					return desc
				}
				return "..."
			}, &ticketKeyWord).
			Options(func() []huh.Option[string] {
				opts := make([]huh.Option[string], len(cfg.TicketKeyWords))
				for i, kw := range cfg.TicketKeyWords {
					opts[i] = huh.NewOption(kw.Name, kw.Name)
				}
				return opts
			}()...).
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
	).WithTheme(theme)

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
	).WithTheme(theme)

	if err := confirmForm.Run(); err != nil {
		fmt.Println("Aborted.")
		return
	}

	if !confirmed {
		fmt.Println("Aborted.")
		return
	}

	cmd := exec.Command("git", "commit", "-m", buildCommitMsg())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", out)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
