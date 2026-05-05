package utils

import "github.com/spf13/viper"

func LoadConfig() (map[string]string, ValidationConfig) {
	// Default values
	commitTypes := DefaultCommitTypes
	validations := ValidationConfig{
		SummaryMaxLen:       60,
		ScopeMaxLen:         30,
		ScopeRequired:       false,
		DescriptionRequired: false,
	}

	// Overwrite with config file if exists
	if viper.IsSet("commit_types") {
		commitTypes = viper.GetStringMapString("commit_types")
	}
	if viper.IsSet("validations.summary_max_length") {
		validations.SummaryMaxLen = viper.GetInt("validations.summary_max_length")
	}
	if viper.IsSet("validations.scope_max_length") {
		validations.ScopeMaxLen = viper.GetInt("validations.scope_max_length")
	}
	if viper.IsSet("validations.require_scope") {
		validations.ScopeRequired = viper.GetBool("validations.require_scope")
	}
	if viper.IsSet("validations.require_description") {
		validations.DescriptionRequired = viper.GetBool("validations.require_description")
	}

	return commitTypes, validations
}
