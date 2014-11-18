package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/path"
	"github.com/bobappleyard/readline"
)

var cdUsage = "cd <path>"

var cd = Command{
	"cd",
	"Changes the current context to the specified path.",
	cdUsage,
	`Changes the current context to the specified path. The path addresses
the name of an object in the project model like "main-server-group".

If the object is part of a collection you can also use an index (zero based)
on the objects type. To avoid naming conflicts you have to prefix the relevant
path segment with ' in that case:

    cd 'hosts[2].'servers[4]

changes the current context to the fifth server of the third host.`,
	// tab completer
	func(project *model.Project, query, _ string) []string {
		readline.CompletionAppendChar = 0
		return path.CurrentPath.Completer(project, query)
	},
	// action
	func(project *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", cdUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", cdUsage)
		}
		return nil
		//		return model.CurrentPath.Cd(project, args[0])
	},
}
