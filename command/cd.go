package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var cdUsage = "cd path"

var cd = Command{
	"cd",
	cdUsage,
	"Changes the current context to the specified path.",
	// tab completer
	func(_, _ string) []string {
		// TODO not yet implemented
		return nil
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", cdUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", cdUsage)
		}

		// TODO validate path
		model.WorkingDir = args[0]
		return nil
	},
}
