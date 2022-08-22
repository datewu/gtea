package utils

import "os"

// InGithubCI check github acitons envrionment
// https://help.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables
func InGithubCI() bool {
	return os.Getenv("GITHUB_ACTION") != ""
}

// PanicFn panic if err is not nil
func PanicFn(fn func() error) {
	err := fn()
	if err != nil {
		panic(err)
	}
}
