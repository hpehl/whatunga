package main

import "strings"

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
	Action func(*Project, []string) error
}

func topLevelCompleter(query, ctx string) []string {
	var results []string
	tokens := strings.Split(ctx, " ")
	if len(tokens) > 0 {
		cmd, valid := commandRegistry[tokens[0]]
		if valid {
			// delegate to command specific completer func
			return cmd.Completer(query, ctx)
		} else {
			for key, _ := range commandRegistry {
				if strings.HasPrefix(key, query) {
					results = append(results, key)
				}
			}
		}
	}
	return results
}

var commandRegistry = make(map[string]Command)
