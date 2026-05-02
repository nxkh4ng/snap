/*
Copyright © 2026 Nguyễn Xuân Khang nxkh4ng

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode/utf8"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	longDesc = `snap is a lightweight CLI tool that helps you make consistent Git Commits without slowing you down.
Following this conventional commits standard - https://www.conventionalcommits.org/en/v1.0.0/`

	typeMap map[string]string
	typeKeys []string
	summaryLen, scopeLen int
	requireScope, requireDescription bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snap",
	Short: "snap your commits into shape",
	Long:  longDesc,
	Run: func(cmd *cobra.Command, args []string) {
		repoCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		if err := repoCmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "not a git repository, use `git init` to initialize")
			os.Exit(1)
		}

		var msg string
		var commit, scope string
		var summary, description string
		var breakingChange string
		var confirm bool

		amend, _ := cmd.Flags().GetBool("amend")
		if amend {
			logCmd := exec.Command("git", "log", "-1", "--format=%B")
			out, err := logCmd.Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to log latest commit: %v\n", err)
				os.Exit(1)
			}
			commit, scope, summary, description, breakingChange, _ = parseCommitMsg(string(out))
		} else {
			autoStage, _ := cmd.Flags().GetBool("all")
			if autoStage {
				addCmd := exec.Command("git", "add", "-u")
				if err := addCmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "failed to stage files: %v\n", err)
					os.Exit(1)
				}
			} else {
				if err := checkGitCommitReady(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		}

		f := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Type*").Value(&commit).
					Placeholder("feat, fix").Suggestions(typeKeys).
					DescriptionFunc(func() string {
						temp := strings.TrimSuffix(strings.TrimSpace(commit), "!")
						if temp == "" {
							return "Select a commit type"
						} else {
							var matchCount int
							var matchedKey string
							for _, key := range typeKeys {
								if strings.HasPrefix(key, temp) {
									matchCount++
									matchedKey = key
									if matchCount > 1 {
										return "Select a commit type"
									}
								}
							}
							if matchCount == 1 {
								return typeMap[matchedKey]
							}
						}
						return "Select a commit type"
					}, &commit).
					Validate(func(t string) error {
						if err := huh.ValidateMinLength(1)(t); err != nil {
							return fmt.Errorf("type cannot be empty")
						}
						t = strings.TrimSuffix(t, "!")
						if _, ok := typeMap[t]; !ok {
							return fmt.Errorf("only allow: %v", strings.Join(typeKeys, ", "))
						}
						return nil
					}),

				huh.NewInput().Value(&scope).
					Placeholder("api, auth").
					TitleFunc(func() string {
						if requireScope {
							return "Scope*"
						}
						return "Scope"
					}, &scope).
					Validate(func(input string) error {
						if err := huh.ValidateMinLength(1)(input); err != nil && requireScope {
							return fmt.Errorf("scope is required")
						}
						if utf8.RuneCountInString(input) > scopeLen {
							current := utf8.RuneCountInString(scope)
							return fmt.Errorf("scope must be at most %d characters long - current: %d", scopeLen, current)
						}
						return nil
					}),

				huh.NewInput().Title("Summary*").Value(&summary).
					Placeholder("Summary of changes").
					Validate(func(input string) error {
						if err := huh.ValidateMinLength(1)(input); err != nil {
							return fmt.Errorf("summary cannot be empty")
						}
						if utf8.RuneCountInString(input) > summaryLen {
							current := utf8.RuneCountInString(summary)
							return fmt.Errorf("summary must be at most %d characters long - current: %d", summaryLen, current)
						}
						return nil
					}),
			),

			huh.NewGroup(
				huh.NewText().Value(&description).
					TitleFunc(func() string {
						if requireDescription {
							return "Description*"
						}
						return "Description"
					}, &description).
					Placeholder("Detailed description of changes").
					Validate(func(input string) error {
						if err := huh.ValidateMinLength(1)(input); err != nil && requireDescription {
							return fmt.Errorf("description is required")
						}
						return nil
					}).
					WithHeight(10),
			),

			huh.NewGroup(
				huh.NewText().Title("BREAKING CHANGE").Value(&breakingChange).
					WithHeight(10),
			).WithHideFunc(func() bool {
				if !strings.Contains(commit, "!") {
					breakingChange = ""
					return true
				}
				return false
			}),
		).WithTheme(huh.ThemeFunc(huh.ThemeBase16))
		if err := f.Run(); err != nil {
			log.Fatal(err)
		}
		msg = formatCommitMsg(commit, scope, summary, description, breakingChange)

		cf := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Value(&confirm).
					TitleFunc(func() string {
						if amend {
							return "Amend this commit?"
						}
						return "Commit this changes?"
					}, &confirm).
					DescriptionFunc(func() string {
						return msg
					}, &msg),
			),
		).WithTheme(huh.ThemeFunc(huh.ThemeBase16))
		if err := cf.Run(); err != nil {
			log.Fatal(err)
		}

		if confirm {
			var commitCmd *exec.Cmd
			if amend {
				commitCmd = exec.Command("git", "commit", "--amend", "-m", msg)
			} else {
				commitCmd = exec.Command("git", "commit", "-m", msg)
			}
			output, err := commitCmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(output))
		}
	},
}

func parseCommitMsg(msg string) (commit, scope, summary, description, breakingChange string, hasBreaking bool) {
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
			commit = before
			scope = strings.TrimSuffix(after, ")")
		} else {
			commit = typeScope
		}

		if hasBreaking {
			commit = commit + "!"
		}
	}

	if before, after, ok := strings.Cut(body, "BREAKING CHANGE:"); ok {
		description = strings.TrimSpace(before)
		breakingChange = strings.TrimSpace(after)
	} else {
		description = body
	}

	return commit, scope, summary, description, breakingChange, hasBreaking
}

func formatCommitMsg(commit, scope, summary, description, breakingChange string) string {
	// cleanup values
	commit = strings.TrimSpace(commit)
	scope = strings.ToLower(strings.TrimSpace(scope))
	summary = strings.TrimSpace(summary)
	description = strings.TrimSpace(description)
	breakingChange = strings.TrimSpace(breakingChange)
	var b strings.Builder

	hasBreakingChange := strings.Contains(commit, "!")
	if hasBreakingChange {
		commit = strings.TrimSuffix(commit, "!")
	}

	// Title = <type>(scope)!: summary
	// or <type>!: summary
	b.WriteString(commit)
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

func checkGitCommitReady() error {
	stagedCmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := stagedCmd.Run()
	if err == nil {
		return fmt.Errorf("no staged changes to commit, use `git add <file>` to stage your changes")
	} else {
		if err.Error() != "exit status 1" {
			return fmt.Errorf("failed to check staged changes: %v", err)
		}
	}
	return nil
}

func loadCommitTypes() map[string]string {
	if viper.IsSet("commit_types") {
		return viper.GetStringMapString("commit_types")
	}
	return map[string]string{
		"feat":     "A new feature",
		"fix":      "A bug fix",
		"docs":     "Documentation only changes",
		"style":    "Formatting, white-space, missing semi-colons,...",
		"refactor": "Code changes that neither fix bugs nor add features",
		"pref":     "Code changes that improves performance",
		"test":     "Adding missing tests or correcting existing tests",
		"build":    "Changes that affect the build system or external dependencies",
		"ci":       "Changes to our CI configuration files and scripts",
		"chore":    "Other changes that don't modify src or test files",
		"revert":   "Reverts a previous commit",
	}
}

func loadValidations() (summaryLen, scopeLen int, requireScope, requireDescription bool) {
	if viper.IsSet("validations.summary_max_length") {
		summaryLen = viper.GetInt("validations.summary_max_length")
	} else {
		summaryLen = 60
	}

	if viper.IsSet("validations.scope_max_length") {
		scopeLen = viper.GetInt("validations.scope_max_length")
	} else {
		scopeLen = 30
	}

	if viper.IsSet("validations.require_scope") {
		requireScope = viper.GetBool("validations.require_scope")
	} else {
		requireScope = false
	}

	if viper.IsSet("validations.require_description") {
		requireDescription = viper.GetBool("validations.require_description")
	} else {
		requireDescription = false
	}

	return summaryLen, scopeLen, requireScope, requireDescription
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snap.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("all", "a", false, "stage all tracked files before committing")
	rootCmd.Flags().Bool("amend", false, "amend the latest commit")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".snap" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".snap")
	}

	viper.AutomaticEnv() // read in environment variables that match

	typeMap = loadCommitTypes()
	typeKeys = make([]string, 0, len(typeMap))
	for key := range typeMap {
		typeKeys = append(typeKeys, key)
	}
	slices.Sort(typeKeys)

	summaryLen, scopeLen, requireScope, requireDescription = loadValidations()
}
