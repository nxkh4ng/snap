/*
Copyright © 2026 Nguyễn Xuân Khang nxkh4ng
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var defaultConfig = `# snap configuration
# https://github.com/nxkh4ng/snap

# Commit types
[commit_types]
feat = "A new feature"
fix = "A bug fix"
docs = "Documentation only changes"
style = "Formatting, white-space, missing semi-colons,..."
refactor = "Code changes that neither fix bugs nor add features"
pref = "Code changes that improves performance"
test = "Adding missing tests or correcting existing tests"
build = "Changes that affect the build system or external dependencies"
ci = "Changes to our CI configuration files and scripts"
chore = "Other changes that don't modify src or test files"
revert = "Reverts a previous commit"

# Validations
[validations]
summary_max_length = 60
scope_max_length = 30
require_scope = false
require_description = false`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize snap config",
	Long:  "Create a `.snap.toml` configuration file in current directory",
	Run: func(cmd *cobra.Command, args []string) {
		filename := ".snap.toml"
		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("File %s already exists\n", filename)
			os.Exit(1)
		}
		err := os.WriteFile(filename, []byte(defaultConfig), 0o644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config file: %v\n", err)
			os.Exit(1)
		}
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get current directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created %s in %s\n", filename, dir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
