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

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	longDesc = `snap is a lightweight CLI tool that helps you make consistent Git Commits without slowing you down.
Following this conventional commits standard - https://www.conventionalcommits.org/en/v1.0.0/`

	types = []string{"feat", "fix", "docs", "chore", "refactor", "test", "revert"}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snap",
	Short: "snap your commits into shape",
	Long:  longDesc,
	Run: func(cmd *cobra.Command, args []string) {
		var commit, scope string
		var summary, description string
		var confirm bool

		f := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Type*").Value(&commit).
					Placeholder("feat, fix").Suggestions(types).
					Validate(func(t string) error {
						t = strings.TrimSpace(t)
						if err := huh.ValidateMinLength(1)(t); err != nil {
							return fmt.Errorf("type cannot be empty")
						}
						if !slices.Contains(types, t) {
							return fmt.Errorf("only allow: %v", strings.Join(types, ", "))
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
				huh.NewText().Title("Description").Value(&description).Placeholder("Detailed description of changes"),
			).WithHeight(9),
		)
		if err := f.Run(); err != nil {
			log.Fatal(err)
		}

		cf := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Title("Commit this changes?").Value(&confirm).
					DescriptionFunc(func() string {
						return formatCommitMsg(commit, scope, summary, description)
					}, confirm),
			),
		)
		if err := cf.Run(); err != nil {
			log.Fatal(err)
		}

		if confirm {
			msg := formatCommitMsg(commit, scope, summary, description)
			fmt.Println(msg)
		}
	},
}

func formatCommitMsg(commit, scope, summary, description string) string {
	// cleanup values
	commit = strings.TrimSpace(commit)
	scope = strings.ToLower(strings.TrimSpace(scope))
	summary = strings.TrimSpace(summary)
	description = strings.TrimSpace(description)

	if scope != "" {
		scope = "(" + scope + ")"
	}
	title := commit + scope + ": " + summary
	var body string
	if description != "" {
		body = "\n\n" + description
	}
	message := title + body
	return message
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
