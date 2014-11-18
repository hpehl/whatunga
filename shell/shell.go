package shell

import (
	"fmt"
	"github.com/bobappleyard/readline"
	"github.com/hpehl/whatunga/command"
	"github.com/hpehl/whatunga/model"
	wpath "github.com/hpehl/whatunga/path"
	"github.com/mitchellh/go-homedir"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	AppName         = "whatunga"
	AppVersionMajor = 0
	AppVersionMinor = 4
	Logo            = `
 __      __.__            __
/  \    /  \  |__ _____ _/  |_ __ __  ____    _________
\   \/\/   /  |  \\__  \\   __\  |  \/    \  / ___\__  \
 \        /|   Y  \/ __ \|  | |  |  /   |  \/ /_/  > __ \_
  \__/\  / |___|  (____  /__| |____/|___|  /\___  (____  /
       \/       \/     \/                \//_____/     \/
`
)

// revision part of the program version.
// This will be set automatically at build time by using:
// go build -ldflags "-X shell.AppVersionRev `date -u +%s`"
var AppVersionRev string

func init() {
	home, err := homedir.Dir()
	if err == nil {
		readline.LoadHistory(path.Join(home, ".whatunga_history"))
	}
}

func lateInit(project *model.Project) {
	readline.Completer = func(query, ctx string) []string {
		// the default which can be overridden by custom command completer functions
		readline.CompletionAppendChar = ' '

		var results []string
		tokens := strings.Fields(ctx)
		//		fmt.Printf("\ntokens: %v\n", tokens)
		if len(tokens) > 0 {
			cmd, valid := command.Registry[tokens[0]]
			if valid {
				// delegate to command completer function
				return cmd.Completer(project, query, ctx)
			} else {
				for key, _ := range command.Registry {
					if strings.HasPrefix(key, query) {
						results = append(results, key)
					}
				}
			}
		} else {
			for key, _ := range command.Registry {
				results = append(results, key)
			}
		}
		return results
	}
}

func Start(info string, project *model.Project) {
	lateInit(project)
	fmt.Printf("%s\n%s\n%s\n", Logo, version(), info)

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
		fmt.Println() // general new line before each cmd execution for better formatting

		tokens := strings.Fields(cmdline)
		cmd, valid := command.Registry[tokens[0]]
		if !valid {
			fmt.Printf("Unknown command: \"%s\"\n", tokens[0])
			prompt(project)
			continue
		}
		if err := cmd.Action(project, tokens[1:]); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}
}

func version() string {
	if len(AppVersionRev) == 0 {
		AppVersionRev = "0"
	}
	return fmt.Sprintf("%s %d.%d.%s (Go runtime %s).",
		AppName, AppVersionMajor, AppVersionMinor, AppVersionRev, runtime.Version())
}

func prompt(project *model.Project) string {
	return fmt.Sprintf("\n[%s:%s @ /%s]> ", project.Name, project.Version, wpath.CurrentPath)
}
