package utils

import "os"

// InGithubCI return whether in github acitons
func InGithubCI() bool {
	// https://help.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables
	return os.Getenv("GITHUB_ACTION") != ""
}
