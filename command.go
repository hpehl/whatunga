package main

import (
	"fmt"
	"os"
)

type Command struct {
	name   string
	action func(*Project) error
}

func nyi(_ *Project) error {
	fmt.Println("Not yet implemented")
	return nil
}

var knownCommands = map[string]Command{
	"show": Command{"show", nyi},
	"exit": Command{"exit", func(_ *Project) error {
		fmt.Println("\nHaere rā")
		os.Exit(0)
		return nil
	}},
}
