package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var setUsage = "set path value,..."

var set = Command{
	"set",
	"Modifies an object / attribute of the project model.",
	setUsage,
	"Modifies an object / attribute of the project model. TODO: Describe path for set.",
	// tab completer
	func(project *model.Project, query, _ string) []string {
		return model.CurrentPath.Completer(project, query)
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing arguments. Usage: %s", setUsage)
		}
		if len(args) > 2 {
			return fmt.Errorf("Too many arguments. Usage: %s", setUsage)
		}
		return nil
	},
}
