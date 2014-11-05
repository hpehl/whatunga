package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	WhatungaJson = "whatunga.json"
	WildFly      = "wildfly"
	EAP          = "eap"
)

type Target struct {
	Name    string
	Version string
}

func (t Target) String() string {
	return t.Name + ":" + t.Version
}

// Used by the flag package to parse a target given as command line flag
func (t *Target) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return errors.New("illegal format.\n")
	}

	var valid = false
	*t = Target{parts[0], parts[1]}
	for _, supported := range supportedTargets {
		if *t == supported {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("unsupported target.\n")
	}

	return nil // no error
}

var supportedTargets = []Target{
	{WildFly, "8.0"},
	{WildFly, "8.1"},
	{EAP, "6.3"},
}

var targetFlag Target = supportedTargets[1]

// default value
var nameFlag string
var versionFlag string

var project *Project

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s [--target=target] [--name=name] [--version=version] <directory>\n", AppName)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Var(&targetFlag, "target", fmt.Sprintf("Specifies the target. Valid targets: %v.", supportedTargets))
	flag.StringVar(&nameFlag, "name", "", "The name of the project. If you omit the name, the directories name is taken.")
	flag.StringVar(&versionFlag, "version", "1.0", `The version which is "1.0" by default.`)
}

func main() {
	flag.Parse()

	// Verify the correct number of arguments
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No directory given!\n\n")
		flag.Usage()
	} else if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Wrong number of arguments!\n\n")
		flag.Usage()
	}

	// if name is not given, use the directory as default
	if nameFlag == "" {
		nameFlag = path.Base(flag.Arg(0))
	}
	// Open existing project or create new one?
	directory := flag.Arg(0)
	fileInfo, err := os.Stat(directory)
	if os.IsNotExist(err) {
		project = newProject(directory, nameFlag, versionFlag, targetFlag)
	} else if fileInfo.Mode().IsDir() {
		// TODO Check for WhatungaJson
		project = openProject(directory)
	} else {
		fmt.Fprintf(os.Stderr, "\"%s\" is not a directory!\n\n", directory)
		flag.Usage()
	}

	// TODO Setup file watcher for whatunga.json and the templates
}

func newProject(directory string, name string, version string, target Target) *Project {
	fmt.Printf("Creating new project %s in directory %s targeting %v....", name, directory, target)

	var perm os.FileMode = 0755

	os.Mkdir(directory, perm)
	os.Chdir(directory)
	os.Mkdir("templates", perm)
	os.Mkdir("downloads", perm)

	fmt.Println("DONE")

	// TODO Create templates

	project = &Project{
		Name:    name,
		Version: version,
		Config: Config{
			Templates: Templates{
				Domain:     "templates/domain.xml",
				HostMaster: "templates/host-master.xml",
				HostSlave:  "templates/host-slave.xml",
			},
			ConsoleUser: User{
				Name:     "admin",
				Password: "passw0rd_",
			},
			DomainUser: User{
				Name:     "domain",
				Password: "passw0rd_",
			},
			DockerRemoteAPI: "unix:///var/run/docker.sock",
		},
		ServerGroups: []ServerGroup{},
		Hosts:        []Host{},
		Users:        []User{},
	}

	b, _ := json.MarshalIndent(project, "", "  ")
	f, _ := os.Create(WhatungaJson)
	defer f.Close()

	f.Write(b)
	f.Sync()

	return project
}

func openProject(directory string) *Project {
	fmt.Printf("Opening existing project in directory %s....", directory)
	os.Chdir(directory)
	fmt.Println("DONE")

	project = &Project{
		Name:    "",
		Version: "",
		Config: Config{
			Templates: Templates{
				Domain:     "templates/domain.xml",
				HostMaster: "templates/host-master.xml",
				HostSlave:  "templates/host-slave.xml",
			},
			ConsoleUser: User{
				Name:     "admin",
				Password: "passw0rd_",
			},
			DomainUser: User{
				Name:     "domain",
				Password: "passw0rd_",
			},
			DockerRemoteAPI: "unix:///var/run/docker.sock",
		},
		ServerGroups: []ServerGroup{},
		Hosts:        []Host{},
		Users:        []User{},
	}
	return project
}
