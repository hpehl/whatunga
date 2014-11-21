package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var rmUsage = "rm <path>"

var rm = Command{
	"rm",
	"Removes an object from the project model.",
	rmUsage,
	"Removes an object from the project model. TODO: Describe path for rm.",
	// tab completer
	func(_ *model.Project, _, _ string) ([]string, int) {
		// TODO not yet implemented
		return nil, 0
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing arguments. Usage: %s", rmUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", rmUsage)
		}
		// TODO not yet implemented
		fmt.Println("Not yet implemented!")
		return nil
	},
}
