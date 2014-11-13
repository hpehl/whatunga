package main

import "fmt"

var rmUsage = "rm path"

var rm = Command{
	"rm",
	rmUsage,
	"Removes an object from the project model.",
	// tab completer
	func(_, _ string) []string {
		// TODO not yet implemented
		return nil
	},
	// action
	func(_ *Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing arguments. Usage: %s", rmUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", rmUsage)
		}
		return nil
	},
}
