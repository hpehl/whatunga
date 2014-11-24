package shell

import (
	"fmt"
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
	AppVersionMinor = 5
	AppVersionMicro = 0
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
	SetWordBreaks(" \t.[:]=")
	home, err := homedir.Dir()
	if err == nil {
		LoadHistory(path.Join(home, ".whatunga_history"))
	}
}

func registerCompleter(project *model.Project) {
	Completer = func(query, ctx string) []string {
		// the default which can be overridden by custom command completer functions
		CompletionAppendChar = ' '

		// The input is exactly one of the commands. In that case we want to add the
		// CompletionAppendChar and go on
		if ctx == query {
			_, exists := command.Registry[query]
			if exists {
				return []string{query}
			}
		}

		var matches []string
		tokens := strings.Fields(ctx)
		if len(tokens) > 0 {
			cmd, valid := command.Registry[tokens[0]]
			if valid {
				// delegate to command completer function
				matches, CompletionAppendChar = cmd.Completer(project, query, ctx)
			} else {
				for key, _ := range command.Registry {
					if strings.HasPrefix(key, query) {
						matches = append(matches, key)
					}
				}
			}
		} else {
			for key, _ := range command.Registry {
				matches = append(matches, key)
			}
		}
		return matches
	}
}

func Start(info string, project *model.Project) {
	registerCompleter(project)
	fmt.Printf("%s\n%s\n%s\n", Logo, version(), info)

	for {
		cmdline, err := Readline(prompt(project))
		if err == io.EOF {
			break // Ctrl-C
		} else if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command: ", err)
			break
		}
		if cmdline == "" {
			continue
		}

		AddHistory(cmdline)
		fmt.Println() // general new line before each cmd execution for better formatting

		tokens := strings.Fields(cmdline)
		cmd, valid := command.Registry[tokens[0]]
		if !valid {
			fmt.Printf("Unknown command: \"%s\"\n", tokens[0])
			continue
		}
		err = cmd.Action(project, tokens[1:])
		if err == command.EXIT {
			AddHistory(cmdline)
			break
		} else if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}

	home, err := homedir.Dir()
	if err == nil {
		SaveHistory(path.Join(home, ".whatunga_history"))
	}
}

func version() string {
	return fmt.Sprintf("%s %d.%d.%d (Go runtime %s).",
		AppName, AppVersionMajor, AppVersionMinor, AppVersionMicro, runtime.Version())
}

func prompt(project *model.Project) string {
	if wpath.CurrentPath.IsEmpty() {
		return fmt.Sprintf("\n[\x1b[0;35m%s\x1b[0;0m:\x1b[0;35m%s\x1b[0;0m] $ ", project.Name, project.Version)
	} else {
		return fmt.Sprintf("\n[\x1b[0;35m%s\x1b[0;0m:\x1b[0;35m%s\x1b[0;0m] \x1b[0;33m%s\x1b[0;0m $ ", project.Name, project.Version, wpath.CurrentPath)
	}
}
