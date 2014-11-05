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

// State constants
type state uint

const (
	COMMANDS      state = iota
	HELP                = iota
	SHOW                = iota
	ADD_TYPE            = iota
	ADD_PATH            = iota
	ADD_VALUE           = iota
	ADD_PARAMETER       = iota
	SET_PATH            = iota
	RM                  = iota
	VALIDATE            = iota
)

//var currentState state = COMMANDS
var workingDir string = "/"

// A function which returns a list of possible options for a state
//type autoComplete func(from state, input string) []string

func shell(info string, project *Project) {
	fmt.Printf("%s\n\n%s\n\n%s\n\n", info, welcome, Version())
	prompt(project)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmdLine := scanner.Text()
		if cmdLine == "" {
			prompt(project)
			continue
		}

		var cmdTokens []string
		words := bufio.NewScanner(strings.NewReader(cmdLine))
		words.Split(bufio.ScanWords)
		for words.Scan() {
			cmdTokens = append(cmdTokens, words.Text())
		}

		cmd, valid := knownCommands[cmdTokens[0]]
		if !valid {
			fmt.Printf("Unknown command: \"%s\"\n", cmdTokens[0])
			prompt(project)
			continue
		}
		cmd.action(project)
		prompt(project)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading command:", err)
	}
}

func prompt(project *Project) {
	fmt.Printf("[%s %s @ %s ]> ", project.Name, project.Version, workingDir)
}
