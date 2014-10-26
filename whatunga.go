package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var project *Project

func init() {
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
}

func main() {
	flag.Parse()
	fmt.Printf("\nArguments: %v", flag.Args())
	fmt.Printf("\nProject: %v", project)

	b, err := json.Marshal(project)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	out.WriteString("\n")
	json.Indent(&out, b, "", "  ")
	out.WriteTo(os.Stdout)

	// TODO Setup file watcher for whatunga.json and the templates
}
