package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hpehl/whatunga/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	WildFly       string      = "wildfly"
	EAP           string      = "eap"
	WhatungaJson  string      = "whatunga.json"
	DirectoryPerm os.FileMode = 0755
	FilePerm      os.FileMode = 0644
)

// ------------------------------------------------------ target

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
	for _, supported := range SupportedTargets {
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

var SupportedTargets = []Target{
	{WildFly, "8.0"},
	{WildFly, "8.1"},
	{EAP, "6.3"},
}

// the xmlns version of the config files
var ModelVersions = map[string]Target{
	"2.0": {WildFly, "8.0"},
	"2.1": {WildFly, "8.1"},
	"1.6": {EAP, "6.3"},
}

// ------------------------------------------------------ project model

type Foo struct {
	Bar Bar `json:"bar"`
}
type Bar struct {
	Name string `json:"name"`
}

type Project struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Config       Config        `json:"config"`
	ServerGroups []ServerGroup `json:"server-groups"`
	Hosts        []Host        `json:"hosts"`
	Users        []User        `json:"users"`
}

func NewProject(directory string, name string, version string, target Target) (*Project, error) {
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
	if err := project.Save(); err != nil {
		return nil, err
	}
	return project, nil
}

func createTemplate(target Target, name string) error {
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

func OpenProject(directory string) (*Project, error) {
	fullyQualifiedWhatungaJson := path.Join(directory, WhatungaJson)
	if _, err := os.Stat(fullyQualifiedWhatungaJson); os.IsNotExist(err) {
		return nil, fmt.Errorf("Missing project file \"%s\"!", fullyQualifiedWhatungaJson)
	}
	if err := os.Chdir(directory); err != nil {
		return nil, err
	}

	var project Project
	if err := project.Load(); err != nil {
		return nil, err
	}
	return &project, nil
}

func (project *Project) Save() error {
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(WhatungaJson, data, FilePerm); err != nil {
		return err
	}
	return nil
}

func (project *Project) Load() error {
	data, err := ioutil.ReadFile(WhatungaJson)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, project); err != nil {
		return err
	}
	return nil
}

func (project *Project) Set(path string, value []string) {
	if path == "" {
		return
	}
	// TODO Use reflection to set values
}

type Config struct {
	Templates       Templates `json:"templates"`
	ConsoleUser     User      `json:"console-user"`
	DomainUser      User      `json:"domain-user"`
	DockerRemoteAPI string    `json:"docker-remote-api"`
}

type Templates struct {
	Domain     string `json:"domain"`
	HostMaster string `json:"host-master"`
	HostSlave  string `json:"host-slave"`
}

type ServerGroup struct {
	Name          string       `json:"name"`
	Profile       string       `json:"profile"`
	SocketBinding string       `json:"socket-binding"`
	Jvm           Jvm          `json:"jvm"`
	Deployments   []Deployment `json:"deployments"`
}

type Deployment struct {
	Name        string `json:"name"`
	RuntimeName string `json:"runtime-name"`
	Path        string `json:"path"`
}

type Host struct {
	Name    string   `json:"name"`
	DC      bool     `json:"domain-controller"`
	Servers []Server `json:"servers"`
	Jvm     Jvm      `json:"jvm"`
}

type Server struct {
	Name        string `json:"name"`
	ServerGroup string `json:"server-group"`
	PortOffset  int    `json:"port-offset"`
	AutoStart   bool   `json:"auto-start"`
	Jvm         Jvm    `json:"jvm"`
}

type BoundedMemory struct {
	Initial string `json:"initial"`
	Max     string `json:"max"`
}

type Jvm struct {
	Name    string        `json:"name"`
	Heap    BoundedMemory `json:"heap"`
	PermGem string        `json:"perm-gen"`
	Stack   string        `json:"stack"`
	Options []string      `json:"options"`
}

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}
