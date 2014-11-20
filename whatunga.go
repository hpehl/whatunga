package main

import (
	"flag"
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/shell"
	"os"
	"path"
)

var targetFlag model.Target = model.SupportedTargets[1]
var nameFlag string
var versionFlag string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--target=target] [--name=name] [--version=version] <directory>\n\n", shell.AppName)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Var(&targetFlag, "target", fmt.Sprintf("Specifies the target. Valid targets: %v.", model.SupportedTargets))
	flag.StringVar(&nameFlag, "name", "", "The name of the project. If you omit the name, the directories name is taken.")
	flag.StringVar(&versionFlag, "version", "1.0", `The project version which is "1.0" by default.`)
}

func main() {
	flag.Parse()

	// Verify the correct number of arguments
	if flag.NArg() == 0 {
		wrongUsage("No directory given!")
	} else if flag.NArg() > 1 {
		wrongUsage("Wrong number of arguments!")
	}

	// if name is not given, use the directory as default
	if nameFlag == "" {
		nameFlag = path.Base(flag.Arg(0))
	}

	// Open existing project or create new one?
	wd, _ := os.Getwd()
	directory := flag.Arg(0)
	fileInfo, err := os.Stat(directory)

	var welcome string
	var project *model.Project
	if os.IsNotExist(err) {
		p, err := model.NewProject(directory, nameFlag, versionFlag, targetFlag)
		if err != nil {
			wrongUsage(err.Error())
		}
		project = p
		welcome = fmt.Sprintf(`Start with new project "%s" in "%s"`, project.Name, path.Join(wd, directory))

	} else if fileInfo.Mode().IsDir() {
		p, err := model.OpenProject(directory)
		if err != nil {
			wrongUsage(err.Error())
		}
		project = p
		welcome = fmt.Sprintf(`Open existing project "%s" in "%s"`, project.Name, path.Join(wd, directory))

	} else {
		wrongUsage(fmt.Sprintf("\"%s\" is not a directory!", directory))
	}

	// TODO Setup file watchers for WhatungaJson and the templates
	shell.Start(welcome, project)
}

func wrongUsage(why string) {
	fmt.Fprintf(os.Stderr, "%s\n\n", why)
	flag.Usage()
}
