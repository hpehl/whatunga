package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
)

var addUsage = "add server-group|host|server|deployment|user <value,...> [--times=n]"

var add = Command{
	"add",
	"Adds one or several objects to the project model.",
	addUsage,
	`Adds one or several objects to the project model. The value(s) depend on the
object type:
    - server-group: The name(s) of the server groups.
    - host:         The name(s) of the hosts.
    - server:       The name(s) of the servers.
    - deployment:   The path to the deployment artifact. Multiple values are not allowed here
    - user:         The username(s) and password(s) separated with ":" as in "foo:bar"

When adding multiple values, you can use a pattern to create unique names.
These patterns can contain specific variables:
    - %w: Resolves to the project name
    - %v: Resolves to the project version
    - %h: Inserts the current host name (applicable when adding servers to a host)
    - %g: Inserts the current server group name (applicable when adding servers to a server group)
    - [n]%c: A counter which starts at zero and which is incremented for each added object.

It's up to the user to choose a pattern which generates unique names.
Non-unique names will lead to an error.`,
	// tab completer
	func(_ *model.Project, _, _ string) ([]string, int) {
		// TODO not yet implemented
		return []string{"not", "yet", "implemented"}, ' '
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing arguments. Usage: %s", addUsage)
		}
		fmt.Println("Not yet implemented!")
		return nil
	},
}
