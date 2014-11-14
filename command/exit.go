package command

import (
	"fmt"
	"github.com/bobappleyard/readline"
	"github.com/hpehl/whatunga/model"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

var exit = Command{
	"exit",
	"Get out of here.",
	"exit",
	"Get out of here.",
	// tab completer
	func(_, _ string) []string {
		return nil
	},
	// action
	func(_ *model.Project, _ []string) error {
		home, err := homedir.Dir()
		if err == nil {
			readline.SaveHistory(path.Join(home, ".whatunga_history"))
		}
		fmt.Println("Haere rā")
		os.Exit(0)
		return nil
	},
}
