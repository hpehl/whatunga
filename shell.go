package main

import (
	"bufio"
	"fmt"
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
	prompt(project)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmdLine := scanner.Text()
		if cmdLine == "" {
			prompt(project)
			continue
		}

		var tokens []string
		words := bufio.NewScanner(strings.NewReader(cmdLine))
		words.Split(bufio.ScanWords)
		for words.Scan() {
			tokens = append(tokens, words.Text())
		}

		cmd, valid := commandRegistry[tokens[0]]
		if !valid {
			fmt.Printf("\nUnknown command: \"%s\"\n", tokens[0])
			prompt(project)
			continue
		}
		if err := cmd(project, tokens[1:]); err != nil {
			fmt.Printf("\n%s\n", err.Error())
		}
		prompt(project)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading command:", err)
	}
}

func prompt(project *Project) {
	fmt.Printf("\n[%s:%s @ %s ]> ", project.Name, project.Version, workingDir)
}
