package command

import (
	"fmt"
	"github.com/hpehl/whatunga/model"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var addSubCommands = []string{"server-group", "host", "server", "deployment", "user"}
var timesOption = "--times"
var timesRegex = regexp.MustCompile("--times=([0-9]+)")
var addUsage = "add " + strings.Join(addSubCommands, "|") + " <value,...> [" + timesOption + "=n]"

var add = Command{
	"add",
	"Adds one or several objects to the project model.",
	addUsage,
	`Adds one or several objects to the project model. The value(s) depend on the
object type:

    - server-group: The name(s) of the server groups.
    - host:         The name(s) of the hosts.
    - server:       The name(s) of the servers.
    - deployment:   The path to the deployment artifact.
                    Multiple values are not allowed here
    - user:         The username(s) and password(s) separated with ":"
                    as in "foo:bar"

To add certain objects, you need to change the context first: To add servers
you need to change the context to a host; to add deployments you need to change
the context to a server group.

When adding multiple values, you can use a pattern to create unique names:

    - %w: Resolves to the project name
    - %v: Resolves to the project version
    - %h: Inserts the current host name (applicable when adding servers
      to a host)
    - %g: Inserts the current server group name (applicable when adding
      deployments to a server group)
    - %[n]c: A counter which starts at zero and which is incremented for
      each added object.

Example:

    cd host[0]
    add server --times=5 %w-%h-server-%1c

Given that there's a host named 'master' and the project's name is 'foo', these
servers are added to the current project model:

    foo-master-server-1
    foo-master-server-2
    foo-master-server-3
    foo-master-server-4
    foo-master-server-5

It's up to the user to choose a pattern which generates unique names.
Non-unique names will lead to an error.`,
	// tab completer
	func(_ *model.Project, query, cmdline string) ([]string, int) {
		var matches []string

		tokens := strings.Fields(cmdline)
		if len(tokens) == 1 && query == "" {
			// just the command was given
			matches = append(matches, addSubCommands...)
			matches = append(matches, timesOption)
		} else {
			var subCommandGiven bool
			var timesOptionGiven bool
			for _, token := range tokens {
				for _, subCommand := range addSubCommands {
					if token == subCommand {
						subCommandGiven = true
						break
					}
				}
				if strings.HasPrefix(token, timesOption) {
					timesOptionGiven = true
				}
			}

			log.SetOutput(os.Stderr)
			log.Printf(`
-+-+-+-+-+
cmdline:          "%s"
query:            "%s"
tokens:           %v
subCommandGiven:  %t
timesOptionGiven: %t
-+-+-+-+-+
`, cmdline, query, tokens, subCommandGiven, timesOptionGiven)

			if !subCommandGiven {
				if query == "" {
					matches = append(matches, addSubCommands...)
				} else {
					for _, subCommand := range addSubCommands {
						if subCommand != query && strings.HasPrefix(subCommand, query) {
							matches = append(matches, subCommand)
						}
					}
				}
			}
			if !timesOptionGiven {
				if strings.HasPrefix(timesOption, query) {
					matches = append(matches, timesOption)
				}
			}
		}

		log.Printf(`
-+-+-+-+-+
matches: %v
-+-+-+-+-+
`, matches)

		if len(matches) == 1 && matches[0] == timesOption {
			return matches, '='
		} else {
			return matches, ' '
		}
	},
	// action
	func(_ *model.Project, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Missing arguments. Usage: %s", addUsage)
		}
		var cmd string
		var values []string
		var times uint64 = 1

		for _, arg := range args {
			if contains(addSubCommands, arg) {
				cmd = arg
			} else if timesRegex.MatchString(arg) {
				groups := timesRegex.FindStringSubmatch(arg)
				if len(groups) < 2 {
					return fmt.Errorf("Error reading \"%s\". Usage: %s", arg, addUsage)
				}
				t, err := strconv.ParseUint(groups[1], 10, 32)
				if err != nil {
					return fmt.Errorf("Error reading \"%s\": %s. Usage: %s", arg, err, addUsage)
				}
				if t < 1 {
					return fmt.Errorf("Invalid value for %s: %d. Usage: %s", timesOption, times, addUsage)
				} else {
					times = t
				}
			} else {
				values = append(values, arg)
			}
		}
		if len(values) == 0 {
			return fmt.Errorf("No values given. Usage: %s", addUsage)
		} else if len(values) > 1 {
			return fmt.Errorf("Too many arguments. Usage: %s", addUsage)
		}

		values = strings.Split(values[0], ",")
		fmt.Printf("Add %d %ss using values %v\n", times, cmd, values)
		return nil
	},
}
