package shell

import (
	"fmt"
	"github.com/bobappleyard/readline"
	"github.com/hpehl/whatunga/command"
	"github.com/hpehl/whatunga/model"
	"github.com/mitchellh/go-homedir"
	"io"
	"os"
	"path"
	"strings"
)

var welcome = `
 __      __.__            __
/  \    /  \  |__ _____ _/  |_ __ __  ____    _________
\   \/\/   /  |  \\__  \\   __\  |  \/    \  / ___\__  \
 \        /|   Y  \/ __ \|  | |  |  /   |  \/ /_/  > __ \_
  \__/\  / |___|  (____  /__| |____/|___|  /\___  (____  /
       \/       \/     \/                \//_____/     \/
`

func init() {
	home, err := homedir.Dir()
	if err == nil {
		readline.LoadHistory(path.Join(home, ".whatunga_history"))
	}
	readline.CompletionAppendChar = ' '
	readline.Completer = func(query, ctx string) []string {
		var results []string
		tokens := strings.Split(ctx, " ")
		if len(tokens) > 0 {
			cmd, valid := command.Registry[tokens[0]]
			if valid {
				// delegate to command specific completer func
				return cmd.Completer(query, ctx)
			} else {
				for key, _ := range command.Registry {
					if strings.HasPrefix(key, query) {
						results = append(results, key)
					}
				}
			}
		}
		return results
	}
}

func Start(info string, project *model.Project) {
	fmt.Printf("%s\n\n%s\n\n%s\n", welcome, Version(), info)

	for {
		cmdline, err := readline.String(prompt(project))
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			break
		}
		if cmdline == "" {
			continue
		}
		readline.AddHistory(cmdline)

		tokens := strings.Split(cmdline, " ")
		var noneEmptyTokens []string
		for _, token := range tokens {
			if strings.TrimSpace(token) != "" {
				noneEmptyTokens = append(noneEmptyTokens, token)
			}
		}
		cmd, valid := command.Registry[tokens[0]]
		if !valid {
			fmt.Printf("\nUnknown command: \"%s\"\n", tokens[0])
			prompt(project)
			continue
		}
		if err := cmd.Action(project, noneEmptyTokens[1:]); err != nil {
			fmt.Printf("\n%s\n", err.Error())
		}
	}
}

func prompt(project *model.Project) string {
	return fmt.Sprintf("\n[%s:%s @ %s ]> ", project.Name, project.Version, model.WorkingDir)
}
