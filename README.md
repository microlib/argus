# Argus

The watchman service. It acts as a deploy/endeploy, start,stop and build of go webservices. It uses the isAlive endpoint in each web service to ensure the webservice is available (if not it will restart it)
The argus service is intsalled using the systemd service ensuring it is always available. A simple script is used to allow for the systemd control 

```bash
systemctl start argus.service

systemctl stop argus.service
``` 
## Outline

Argus is made up of a client and server component

### Server
* Deploy - allows the cloning of a git repo to a directory on the server
* Undeploy - removes the repo from the server
* Start - executes a script (see Note below) to start the go web service
* Stop - executes a script to stop the web service
* Build - builds the go executable

A systemd service is use to ensure argus is always available (restart policy of 5 secs after a crash)
Argus uses the 'isalive' endpoint on each web service to ensure high availability and will restart the service if its down
The obvious choice of using a simple start/stop script is due to permissions (the rest api interface) does not have sudo or su permissions and
therefore can't install a systemd service script.

### Client
Interfaces to the application-server list (via a config.json file)
It will try deploy to the server providing the server has enough resources (via getserverstats endpoint)
If the server does not have enough resources it will try the next server in the list unitl the complete server list has been checked
The client has the functionaliy to deploy,undeploy,start,stop and build a webservice (golang only)

Example usage (to deploy) :

```bash

./argus-client 21423432 go-simple-service deploy deploy https://github.com/dimitraz

```

Example usage (to start/stop) :

```bash

./argus-client 21423432 go-simple-service appexecute start https://github.com/dimitraz

./argus-client 21423432 go-simple-service appexecute stop https://github.com/dimitraz

```

## Installation

TBD - Ansible will be used to install argus. It will use an inventory file listing the available "application-servers" to install to
SSH access and sudo permissions will be needed to install argus as it makes use of systemd

## Issues

This is still a WIP and as this is an initial POC here is a list of known issues :-
* Update HA-Proxy config (virtual server) for proxy and reverse proxy of each deployed web service
* Token (will make use of JWT)
* Authentication and authorization - integrate the acl project here 
* Golang only 
* Logging
* No unit testing (bad I know - TDD laziness)


## Note
The http server by @luigizuccarelli uses signals to allow for graceful shutdown. Use this as a standard pattern when creating all web services. 

