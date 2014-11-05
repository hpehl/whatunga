package main

import (
	"bufio"
	"fmt"
	"os"
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

var currentState state = COMMANDS
var workingDir string = "/"

// A function which returns a list of possible options for a state
type autoComplete func(from state, input string) []string

func shell(info string, project *Project) {
	fmt.Printf("%s\n\n%s\n\n%s\n\n", info, welcome, Version())
	scanner := bufio.NewScanner(os.Stdin)
	prompt(project)
	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd == "exit" {
			break
		}
		fmt.Println(cmd)
		prompt(project)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func prompt(project *Project) {
	fmt.Printf("[%s %s @ %s ]> ", project.Name, project.Version, workingDir)
}
