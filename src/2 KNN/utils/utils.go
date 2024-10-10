package utils

import "flag"

// https://www.educative.io/answers/how-to-check-if-a-command-line-flag-is-set-in-go
func IsFlagPassed(flagName string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			found = true
		}
	})
	return found
}
