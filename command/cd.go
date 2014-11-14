package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var cdUsage = "cd <path>"

var cd = Command{
	"cd",
	"Changes the current context to the specified path.",
	cdUsage,
	`Changes the current context to the specified path. The path addresses an
object or attribute in the project model. Could be a specific attribute like
"host-master.server1.port-offset" or an object like "main-server-group".

If the object is part of a collection you can also use an index (zero based)
on objects type. To avoid naming conflicts you have to prefix the relevant
path segment with ' in that case:

    set 'hosts[2].servers[4].auto-start true

sets the auto start flag of the fifth server of the third host`,
	// tab completer
	func(_ *model.Project, _, _ string) []string {
		// TODO not yet implemented
		return nil
	},
	// action
	func(project *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing argument. Usage: %s", cdUsage)
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", cdUsage)
		}

		path := model.Path(args[0])
		if err := path.Validate(project); err != nil {
			return err
		}
		model.CurrentContext = path
		return nil
	},
}
