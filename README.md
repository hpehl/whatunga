# Usage 

You can use whatunga to create a new project or to open an existing one. The general syntax is 

	whatunga [--target=target] [--name=name] [--version=version] <directory>
	
Whatunga looks for a file named `whatunga.json` in the specified directory. If there's one, whatunga opens the related project. Otherwise a new empty project is created in the given directory. 

## Options

The options are only valid when creating a new project; they're ignored when an existing project is loaded. 

- `target` Specifies the target using `<product>:<version>` with the following valid combinations:

	- wildfly:8.0
	- wildfly:8.1
	- eap:6.3
	
	If you omit the target, `wildfly:8.1` will be used.
	
- `name` The name of the project. If you omit the name, the directories name is taken.

- `version` The version which is "1.0" by default.

# Model

Whatunga stores all configuration, server groups, hosts, servers, deployments and other settings in a JSON file called `whatunga.json`. You can also edit this file externally. Whatunga will watch the file for changes and reload its internal state whenever the file is changed. Roughly the JSON file consists of these sections:

- name & version
- configuration
- domain model (server groups, hosts, servers, deployments and other domain settings)
- users

The following example shows a typical `whatunga.json` file:

```
{
  "name": "eq08",
  "version": "0.5",
  "config": {
    "templates": {
      "domain": "templates/domain.xml",
      "HostMaster": "templates/host-master.xml",
      "HostSlave": "templates/host-slave.xml"
    },
    "console-user": {
      "username": "admin",
      "password": "passw0rd_"
    },
    "domain-user": {
      "username": "domain",
      "password": "passw0rd_"
    },
    "docker-remote-api": "unix:///var/run/docker.sock"
  },
  "server-groups": [
    {
      "name": "default-group",
      "profile": "default-profile",
      "socket-binding": "default-sockets",
      "deployments": [
        {
          "name": "ticketmonster.ear",
          "runtime-name": "ticketmonster",
          "path": "deployments/ticketmonster.ear"
        }
      ]
    },
    {
      "name": "full-group",
      "profile": "full-profile",
      "socket-binding": "full-sockets",
      "jvm": {
        "name": "default-jvm",
        "heap": {
          "initial": "1GB",
          "max": "2GB"
        },
        "perm-gen": "256MB",
        "stack": "128MB"
      }
    },
    {
      "name": "ha-group",
      "profile": "ha-profile",
      "socket-binding": "ha-sockets"
    }
  ],
  "hosts": [
    {
      "name": "master",
      "domain-controller": true,
      "servers": [
        {
          "name": "server-one",
          "server-group": "default-group"
        },
        {
          "name": "server-two",
          "server-group": "default-group",
          "port-offset": 50
        },
        {
          "name": "server-three",
          "server-group": "full-group",
          "port-offset": 100,
          "auto-start": true
        }
      ]
    },
    {
      "name": "slave1",
      "servers": [
        {
          "name": "server-four",
          "server-group": "full-group"
        },
        {
          "name": "server-five",
          "server-group": "full-group",
          "port-offset": 50
        }
      ]
    },
    {
      "name": "slave2",
      "servers": [
        {
          "name": "server-six",
          "server-group": "ha-group"
        },
        {
          "name": "server-seven",
          "server-group": "ha-group",
          "port-offset": 50,
          "auto-start": true
        }
      ]
    }
  ],
  "users": [
    {
      "username": "monitor",
      "password": "passw0rd_"
    },
    {
      "username": "test",
      "password": "passw0rd_"
    }
  ]
}
```

In the following sections the basic parts of the project model are described in more detail.

## Configuration

Whatunga uses a set of configuration properties and default values. You can change these values using the commands below or by editing `whatungs.json` directly. The configuration consists of these settings

- templates
- mandatory users
- docker endpoint

### Templates

Templates are used to generate the final configuration files. The domain template defines the available profiles and socket-bindings which can be used when adding server groups, hosts and servers. 

Per default the templates are located in a folder called `templates` relative to the main configuration file. For new projects they will be generated based on the chosen WildFly / EAP version. You can change the path anytime using:

	set config.templates.domain /your/path/for/domain.xml
	
Please make sure that you don't mix templates from different products and versions. A domain template for WildFly 8.1 won't work with a host-master template targeting EAP 6.3. 
	
### User

The configuration includes two fixed users, which must not be removed: 

1. A user for the management interfaces (CLI / Admin Console). 
1. A user for the connection between the domain controller and the slaves.

Both users are added to the docker containers using the `add-user` script. 

### Docker

In order to generate and start the WildFly / EAP instances the remote Docker API is used. The endpoint is stored under the configuration property `config.docker-remote-api`.

## Domain Settings

These settings hold the actual domain model. Server groups, hosts, servers and deployments are stored here. Use the commands described below to add additional objects.  

### Deployments

Deployments artifacts can live anywhere on the filesystem and you can add them using their absolute filename. However it is considered as good practice to copy the deployment artifacts to the folder `deployments` - relative to the project - before adding them to the project model. Doing so will make sure you'll end up with a self containing project with all necessary resource relative to the project root. 

In order to prevent naming problems when using paths (see below), the deployment name is based on the file name, but points are replaced with dashes. 

## Users

In this section you can add additional users which are added to the domain controller using the `add-user` script.

# Commands

Whatunga provides a list of commands to show current settings, change the project model and interact with Docker.

- `help [command]` Shows the list of available commands or context sensitive help

- `cd path` Changes the current context to the specified path.

- `ls [path]` Lists the model of the current context or specified path.

- `add server-group|host|server|deployment|user value,... [--times=n]` Adds one or several objects to the project model.

- `set path value,...` Modifies an object / attribute of the project model.

- `rm path` Removes an object from the project model.

- `validate` Checks whether the project model is valid.

- `docker cmd` Docker related commands
	- `create` Creates docker images based on the current project model.
	- `start` Starts the docker images.

- `exit` Get out of here.

## Path

For some commands you need to provide a path. A path specifies an object or attribute in the project model. This could be a specific attribute like `hosts[host-master].servers[0].port-offset` or an object like `server-groups[main-server-group]`. For bulk operations the path can include a range (which follows the [rules for slices](http://tour.golang.org/#33) in the Go language).

To set the auto start flag of all servers in the group `staging-group` use the following command:

    set server-groups[staging-group].servers[:].auto-start true
    
## Value

Values can be simple values like `100`, `true` or `"128MB"` or full JSON encoded objects like `{"name":"s2jvm","heap":{"initial":"1GB","max":"2GB"},"options":["-server"]}`.

If the path specified multiple objects you can provide multiple values:

	# Set the auto start flag of the servers of host master to the given values
	set hosts[master].servers[:].auto-start true,false,false,true

## Naming

When adding multiple server groups, hosts and servers, whatunga uses a naming pattern to create unique names. These patterns can contain specific variables: 

- `%w` Resolves to the project name
- `%v` Resolves to the project version
- `%h` Inserts the current host name (applicable when adding servers to a host)
- `%g` Inserts the current server group name (applicable when adding servers to a server group)
- `[n]%c` A counter which starts at zero and which is incremented for each added object.

It's up to the user to choose a pattern which generates unique names. Non-unique names will lead to an error.  

## Examples

The following sample shows a list of commands to setup a domain with three server groups, four hosts, eight servers and one deployment:

```
# Add three server groups and set the specified profiles. 
# As no socket binding is specified, the first socket binding defined 
# in the domain template is used. 
add server-group group%c --times=3
set server-groups[:].profile dev,staging,prod
    
# Assign a deployment to the dev server group
cd server-groups[0]
add deployment deployments/ticketmonster.ear
set deployments[0].runtime-name ticketmonster
cd ..

# Add a single host named "master" using default values. After that the
# domain controller flag is set to true. Finally add three more hosts.
add host master
set hosts[master].domain-controller true
add host slave0,slave1,slave2

# For each slave add servers. By default the port offset is incremented
# by 50,  the server group is assigned to the first defined server group 
# and auto start is disabled 
cd hosts[slave0]
add server server%c --times=3
cd ..
cd hosts[slave1]
add server server%3c --times=3
cd ..
cd hosts[slave2]
add server server6,lastserver
cd ..

# Change defaults
set hosts[slave0].servers[server0].auto-start true

set hosts[slave1].servers[0:3].server-group staging,staging,qa
set hosts[slave1].servers[0:3].auto-start false,true,true

set hosts[slave2].servers[:].server-group qa,qa,prod
set hosts[slave2].servers[:].port-offset 0+20

cd hosts[2].servers[2]
set jvm {"name":"s2jvm","heap":{"initial":"1GB","max":"2GB"},"options":["-server"]}
cd /
```

# Limitations

Users are stored in the properties based user store. The usage of an external user store like LDAP / ActiveDirectory is not yet supported.