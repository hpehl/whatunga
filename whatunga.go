package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hpehl/whatunga/model"
	"github.com/hpehl/whatunga/shell"
	"github.com/hpehl/whatunga/template"
	"io/ioutil"
	"os"
	"path"
)

const (
	WhatungaJson  string      = "whatunga.json"
	DirectoryPerm os.FileMode = 0755
	FilePerm      os.FileMode = 0644
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

	if os.IsNotExist(err) {
		project, err := newProject(directory, nameFlag, versionFlag, targetFlag)
		if err != nil {
			wrongUsage(err.Error())
		}
		shell.Start(fmt.Sprintf(`Start with new project "%s" in "%s"`, project.Name, path.Join(wd, directory)), project)

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
		shell.Start(fmt.Sprintf(`Open existing project "%s" in "%s"`, project.Name, path.Join(wd, directory)), project)

	} else {
		wrongUsage(fmt.Sprintf("\"%s\" is not a directory!", directory))
	}

	// TODO Setup file watchers for WhatungaJson and the templates
}

func wrongUsage(why string) {
	fmt.Fprintf(os.Stderr, "%s\n\n", why)
	flag.Usage()
}

func newProject(directory string, name string, version string, target model.Target) (*model.Project, error) {
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

	project := &model.Project{
		Name:    name,
		Version: version,
		Config: model.Config{
			Templates: model.Templates{
				Domain:     "templates/domain.xml",
				HostMaster: "templates/host-master.xml",
				HostSlave:  "templates/host-slave.xml",
			},
			ConsoleUser: model.User{
				Name:     "admin",
				Password: "passw0rd_",
			},
			DomainUser: model.User{
				Name:     "dc",
				Password: "passw0rd_",
			},
			DockerRemoteAPI: "unix:///var/run/docker.sock",
		},
		ServerGroups: []model.ServerGroup{},
		Hosts:        []model.Host{},
		Users:        []model.User{},
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

func createTemplate(target model.Target, name string) error {
	templatePath := path.Join("templates", target.Name, target.Version, name)
	data, err := template.Asset(templatePath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join("templates", name), data, FilePerm); err != nil {
		return err
	}
	return nil
}

func openProject(directory string) (*model.Project, error) {
	os.Chdir(directory)

	var project model.Project
	data, err := ioutil.ReadFile(WhatungaJson)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}
	return &project, nil
}
