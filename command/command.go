package command

import "github.com/hpehl/whatunga/model"

var cmdOrder []string

type registry map[string]Command

func (r registry) Add(cmd Command) {
	r[cmd.Name] = cmd
	cmdOrder = append(cmdOrder, cmd.Name)
}

// use internal iteration to control the iteration order
func (r registry) forEach(fn func(cmd Command)) {
	for _, name := range cmdOrder {
		fn(r[name])
	}
}

var Registry = make(registry)

func init() {
	Registry.Add(show)
	Registry.Add(cd)
	Registry.Add(add)
	Registry.Add(set)
	Registry.Add(rm)
	Registry.Add(validate)
	Registry.Add(docker)
	Registry.Add(exit)
	Registry.Add(help)
}

type Command struct {
	// The name of the command
	Name string
	// A short description of this command
	Description string
	// How to use this command
	Usage string
	// A more detailed description for this command
	UsageDescription string
	// The function to call when checking for tab completion
	Completer func(query, ctx string) []string
	// The function to call when this command is invoked
	Action func(*model.Project, []string) error
}
