package main

import (
	"encoding/json"
	"fmt"
)

type Project struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Config       Config        `json:"config"`
	ServerGroups []ServerGroup `json:"server-groups"`
	Hosts        []Host        `json:"hosts"`
	Users        []User        `json:"users"`
}

func (project *Project) Set(path string, value []string) {
	if path == "" {
		return
	}

}

type Config struct {
	Templates       Templates `json:"templates"`
	ConsoleUser     User      `json:"console-user"`
	DomainUser      User      `json:"domain-user"`
	DockerRemoteAPI string    `json:"docker-remote-api"`
}

func (c Config) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error generating configuration: %s", err.Error())
	}
	return string(data)
}

type Templates struct {
	Domain     string `json:"domain"`
	HostMaster string `json:"host-master`
	HostSlave  string `json:host-slave`
}

type ServerGroup struct {
	Name          string       `json:"name"`
	Profile       string       `json:"profile"`
	SocketBinding string       `json:"socket-binding"`
	Jvm           Jvm          `json:"jvm"`
	Deployments   []Deployment `json:"deployments"`
}

func (sg ServerGroup) String() string {
	data, err := json.MarshalIndent(sg, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error generating server group: %s", err.Error())
	}
	return string(data)
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

func (h Host) String() string {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error generating host: %s", err.Error())
	}
	return string(data)
}

type Server struct {
	Name        string `json:"name"`
	ServerGroup string `json:"server-group"`
	PortOffset  uint   `json:"port-offset"`
	AutoStart   bool   `json:"auto-start"`
	Jvm         Jvm    `json:"jvm"`
}

type Memory string

type BoundedMemory struct {
	Initial Memory `json:"initial"`
	Max     Memory `json:"max"`
}

type Jvm struct {
	Name    string        `json:"name"`
	Heap    BoundedMemory `json:"heap"`
	PermGem Memory        `json:"perm-gen"`
	Stack   Memory        `json:"stack"`
	Options []string      `json:"options"`
}

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}
