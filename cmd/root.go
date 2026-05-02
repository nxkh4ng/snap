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
	"slices"
	"strings"
	"unicode/utf8"

	"charm.land/huh/v2"
	"github.com/nxkh4ng/snap/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	longDesc = `snap is a lightweight CLI tool that helps you make consistent Git Commits without slowing you down.
Following this conventional commits standard - https://www.conventionalcommits.org/en/v1.0.0/`

	typeMap map[string]string
	validations internal.ValidationConfig
	typeKeys []string
	summaryMaxLen, scopeMaxLen int
	scopeRequired, descriptionRequired bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snap",
	Short: "snap your commits into shape",
	Long:  longDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.CheckGitRepo(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var msg string
		var commit, scope string
		var summary, description string
		var breakingChange string
		var confirm bool

		amendFlag, _ := cmd.Flags().GetBool("amend")
		if amendFlag {
			latestMsg, err := internal.GetTheLatestCommitMsg()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			commit, scope, summary, description, breakingChange, _ = internal.ParseCommitMsg(latestMsg)
		} else {
			autoStageFlag, _ := cmd.Flags().GetBool("all")
			if autoStageFlag {
				if err := internal.StageAll(); err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					os.Exit(1)
				}
			} else {
				if err := internal.CheckGitCommitReady(); err != nil {
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
						if scopeRequired {
							return "Scope*"
						}
						return "Scope"
					}, &scope).
					Validate(func(input string) error {
						if err := huh.ValidateMinLength(1)(input); err != nil && scopeRequired {
							return fmt.Errorf("scope is required")
						}
						if utf8.RuneCountInString(input) > scopeMaxLen {
							current := utf8.RuneCountInString(scope)
							return fmt.Errorf("scope must be at most %d characters long - current: %d", scopeMaxLen, current)
						}
						return nil
					}),

				huh.NewInput().Title("Summary*").Value(&summary).
					Placeholder("Summary of changes").
					Validate(func(input string) error {
						if err := huh.ValidateMinLength(1)(input); err != nil {
							return fmt.Errorf("summary cannot be empty")
						}
						if utf8.RuneCountInString(input) > summaryMaxLen {
							current := utf8.RuneCountInString(summary)
							return fmt.Errorf("summary must be at most %d characters long - current: %d", summaryMaxLen, current)
						}
						return nil
					}),
			),

			huh.NewGroup(
				huh.NewText().Value(&description).
					TitleFunc(func() string {
						if descriptionRequired {
							return "Description*"
						}
						return "Description"
					}, &description).
					Placeholder("Detailed description of changes").
					Validate(func(input string) error {
						if err := huh.ValidateMinLength(1)(input); err != nil && descriptionRequired {
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
		msg = internal.FormatCommitMsg(commit, scope, summary, description, breakingChange)

		cf := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Value(&confirm).
					TitleFunc(func() string {
						if amendFlag {
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
			output, err := internal.Commit(msg, amendFlag)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(output)
		}
	},
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

	typeMap, validations = internal.LoadConfig()
	typeKeys = make([]string, 0, len(typeMap))
	for key := range typeMap {
		typeKeys = append(typeKeys, key)
	}
	slices.Sort(typeKeys)

	summaryMaxLen = validations.SummaryMaxLen
	scopeMaxLen = validations.ScopeMaxLen
	scopeRequired = validations.ScopeRequired
	descriptionRequired = validations.DescriptionRequired
}
