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

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	longDesc = `snap is a lightweight CLI tool that helps you make consistent Git Commits without slowing you down.
Following this conventional commits standard - https://www.conventionalcommits.org/en/v1.0.0/`

	typeMap = map[string]string{
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

	typeKeys []string
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
		var msg string
		var commit, scope string
		var summary, description string
		var breakingChange string
		var confirm bool

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

				huh.NewInput().Title("Scope").Value(&scope).
					Placeholder("api, auth").CharLimit(30),

				huh.NewInput().Title("Summary*").Value(&summary).
					Placeholder("Summary of changes").CharLimit(60).
					Validate(func(input string) error {
						if err := huh.ValidateMinLength(1)(input); err != nil {
							return fmt.Errorf("summary cannot be empty")
						}
						return nil
					}),
			),

			huh.NewGroup(
				huh.NewText().Title("Description").Value(&description).
					Placeholder("Detailed description of changes").
					WithHeight(10),
			),

			huh.NewGroup(
				huh.NewText().Title("BREAKING CHANGE").Value(&breakingChange).
					WithHeight(10),
			).WithHideFunc(func() bool {
				return !strings.Contains(commit, "!")
			}),
		)
		if err := f.Run(); err != nil {
			log.Fatal(err)
		}
		msg = formatCommitMsg(commit, scope, summary, description, breakingChange)

		cf := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Title("Commit this changes?").Value(&confirm).
					DescriptionFunc(func() string {
						return msg
					}, &confirm),
			),
		)
		if err := cf.Run(); err != nil {
			log.Fatal(err)
		}

		if confirm {
			commitCmd := exec.Command("git", "commit", "-m", msg)
			output, err := commitCmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(output))
		}
	},
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	typeKeys = make([]string, 0, len(typeMap))
	for key := range typeMap {
		typeKeys = append(typeKeys, key)
	}
	slices.Sort(typeKeys)

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snap.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("all", "a", false, "stage all tracked files before committing")
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
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".snap")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
