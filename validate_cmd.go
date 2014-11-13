package main

import "fmt"

var validateUsage = "validate"

var validate = Command{
	"validate",
	validateUsage,
	"Checks whether the project model is valid.",
	// tab completer
	func(_, _ string) []string {
		return nil
	},
	// action
	func(_ *Project, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Illegal argument. Usage: %s", validateUsage)
		}
		return nil
	},
}
