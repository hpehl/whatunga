package command

import (
	"errors"
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var EXIT = errors.New("exit")

var exit = Command{
	"exit",
	"Get out of here.",
	"exit",
	"Get out of here.",
	// tab completer
	func(_ *model.Project, _, _ string) ([]string, int) {
		return nil, 0
	},
	// action
	func(_ *model.Project, _ []string) error {
		fmt.Println("Haere rā\n")
		return EXIT
	},
}
