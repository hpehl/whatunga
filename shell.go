package main

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

// A function which returns a list of possible options for a state
type autoComplete func(st state, input string) []string
