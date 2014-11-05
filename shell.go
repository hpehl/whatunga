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
var WorkingDir string = "/"

// A function which returns a list of possible options for a state
//type autoComplete func(from state, input string) []string

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

		cmd, valid := CommandRegistry[tokens[0]]
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
	fmt.Printf("\n[%s:%s @ %s ]> ", project.Name, project.Version, WorkingDir)
}
