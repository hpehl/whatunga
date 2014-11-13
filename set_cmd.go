package main

import "fmt"

var setUsage = "set path value,..."

var set = Command{
	"set",
	setUsage,
	"Modifies an object / attribute of the project model.",
	// tab completer
	func(_, _ string) []string {
		// TODO not yet implemented
		return nil
	},
	// action
	func(_ *Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing arguments. Usage: %s", setUsage)
		}
		if len(args) > 2 {
			return fmt.Errorf("Too many arguments. Usage: %s", setUsage)
		}
		return nil
	},
}
