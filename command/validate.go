package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var validateUsage = "validate"

var validate = Command{
	"validate",
	"Checks whether the project model is valid.",
	validateUsage,
	"Checks whether the project model is valid.",
	// tab completer
	func(_ *model.Project, _, _ string) ([]string, int) {
		return nil, 0
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("Illegal argument. Usage: %s", validateUsage)
		}
		// TODO not yet implemented
		fmt.Println("Not yet implemented!")
		return nil
	},
}
