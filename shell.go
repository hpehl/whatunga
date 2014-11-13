package main

import (
	"fmt"
	"github.com/bobappleyard/readline"
	"io"
	"os"
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

//var currentState state = COMMANDS
var workingDir string = "/"

func shell(info string, project *Project) {
	fmt.Printf("%s\n\n%s\n\n%s\n", info, welcome, Version())

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
		cmd, valid := commandRegistry[tokens[0]]
		if !valid {
			fmt.Printf("\nUnknown command: \"%s\"\n", tokens[0])
			prompt(project)
			continue
		}
		if err := cmd.Action(project, tokens[1:]); err != nil {
			fmt.Printf("\n%s\n", err.Error())
		}
	}
}

func prompt(project *Project) string {
	return fmt.Sprintf("\n[%s:%s @ %s ]> ", project.Name, project.Version, workingDir)
}
