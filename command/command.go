package command

import "github.com/hpehl/whatunga/model"

var Registry = make(map[string]Command)

func init() {
	// init commands
	Registry["help"] = help
	Registry["show"] = show
	Registry["cd"] = cd
	Registry["add"] = add
	Registry["set"] = set
	Registry["rm"] = rm
	Registry["validate"] = validate
	Registry["docker"] = docker
	Registry["exit"] = exit
}

type Command struct {
	// The name of the command
	Name string
	// A short description of the usage of this command
	Usage string
	// Help text
	Help string
	// The function to call when checking for bash command completions
	Completer func(query, ctx string) []string
	// The function to call when this command is invoked
	Action func(*model.Project, []string) error
}
