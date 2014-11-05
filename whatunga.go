package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	WhatungaJson  string      = "whatunga.json"
	WildFly       string      = "wildfly"
	EAP           string      = "eap"
	DirectoryPerm os.FileMode = 0755
	FilePerm      os.FileMode = 0644
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
var nameFlag string
var versionFlag string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--target=target] [--name=name] [--version=version] <directory>\n\n", AppName)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Var(&targetFlag, "target", fmt.Sprintf("Specifies the target. Valid targets: %v.", supportedTargets))
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
	directory := flag.Arg(0)
	fileInfo, err := os.Stat(directory)

	if os.IsNotExist(err) {
		project, err := newProject(directory, nameFlag, versionFlag, targetFlag)
		if err != nil {
			wrongUsage(err.Error())
		}
		shell(fmt.Sprintf(`Start with new project "%s" in "%s"`, project.Name, directory), project)

	} else if fileInfo.Mode().IsDir() {
		// Check for WhatungaJson
		fullyQualifiedWhatungaJson := path.Join(directory, WhatungaJson)
		if _, err := os.Stat(fullyQualifiedWhatungaJson); os.IsNotExist(err) {
			wrongUsage(fmt.Sprintf("Missing project file \"%s\"!", fullyQualifiedWhatungaJson))
		}
		project, err := openProject(directory)
		if err != nil {
			wrongUsage(err.Error())
		}
		shell(fmt.Sprintf(`Open existing project "%s" in "%s"`, project.Name, directory), project)

	} else {
		wrongUsage(fmt.Sprintf("\"%s\" is not a directory!", directory))
	}

	// TODO Setup file watchers for WhatungaJson and the templates
}

func wrongUsage(why string) {
	fmt.Fprintf(os.Stderr, "%s\n\n", why)
	flag.Usage()
}

func newProject(directory string, name string, version string, target Target) (*Project, error) {
	if err := os.MkdirAll(directory, DirectoryPerm); err != nil {
		return nil, err
	}
	if err := os.Chdir(directory); err != nil {
		return nil, err
	}
	if err := os.Mkdir("templates", DirectoryPerm); err != nil {
		return nil, err
	}
	if err := createTemplate(target, "domain.xml"); err != nil {
		return nil, err
	}
	if err := createTemplate(target, "host-master.xml"); err != nil {
		return nil, err
	}
	if err := createTemplate(target, "host-slave.xml"); err != nil {
		return nil, err
	}
	if err := os.Mkdir("downloads", DirectoryPerm); err != nil {
		return nil, err
	}

	project := &Project{
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
				Name:     "dc",
				Password: "passw0rd_",
			},
			DockerRemoteAPI: "unix:///var/run/docker.sock",
		},
		ServerGroups: []ServerGroup{},
		Hosts:        []Host{},
		Users:        []User{},
	}

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(WhatungaJson, data, FilePerm); err != nil {
		return nil, err
	}

	return project, nil
}

func createTemplate(target Target, name string) error {
	template := path.Join("templates", target.Name, target.Version, name)
	data, err := Asset(template)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join("templates", name), data, FilePerm); err != nil {
		return err
	}
	return nil
}

func openProject(directory string) (*Project, error) {
	os.Chdir(directory)

	var project Project
	data, err := ioutil.ReadFile(WhatungaJson)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}

	return &project, nil
}
