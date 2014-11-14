package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var addUsage = "add server-group|host|server|deployment|user value,... [--times=n]"

var add = Command{
	"add",
	addUsage,
	"Adds one or several objects to the project model.",
	// tab completer
	func(_, _ string) []string {
		// TODO not yet implemented
		return nil
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing arguments. Usage: %s", addUsage)
		}
		return nil
	},
}
