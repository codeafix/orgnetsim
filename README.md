# <img src="https://github.com/codeafix/orgnetsim/blob/main/web/src/logo.svg" width="60" height="60" /> orgnetsim [![Build Status](https://github.com/codeafix/orgnetsim/actions/workflows/Build.yml/badge.svg?branch=main)](https://github.com/codeafix/orgnetsim/actions/workflows/Build.yml) [![codecov](https://codecov.io/github/codeafix/orgnetsim/graph/badge.svg?token=nMO0zwQVEY)](https://codecov.io/github/codeafix/orgnetsim) [![MIT](https://img.shields.io/npm/l/express.svg)](https://github.com/codeafix/orgnetsim/blob/main/LICENSE)
A simulator for Organisational Networks

The simulator is created from a Network of Agents. The Network itself can be any arbitrary graph and contains a collection of Agents and a collection of links between those Agents. The simulator uses Colors to represent competing ideas on the Network. The default Color for an Agent is Grey. During a simulation Agents interact and decide whether or not to update their Color.

## Packages

[sim](sim/README.md) The organisation network simulator

[srvr](srvr/README.md) The organisation network simulator web server

[web](web/README.md) A REACT based front-end for the simulator

## Command line utility

[orgnetsim](orgnetsim/README.md) A command line utility for parsing lists and creating networks

## Docker
This project also contains a Dockerfile that builds both the API and front-end into a container. All data that is created in the app is stored in the container path `/tmp/data`. By default the container will be built with an empty simulation list. To persist data outside the container make sure to mount the container path to a persistable storage path on the host machine, or build and run using the ```docker.compose``` which sets this up.
```
docker run -v <host_path>:/tmp/data -d -p 8080:8080 orgnetsim:v0.1
```
The docker file is set up to run as a non-root user. In order for this to work properly the container must be set up with a user that has sufficient permission to read and write files in the persistable storage path on the host machine. There are two arguments ```UID``` and ```GID``` to the Dockerfile that can be used to set user id and group id of the user that is created in the container by the docker file. Make sure these are set to a user that exists on the host machine that has read and write permissions to the persistable storage path. The docker file creates a group in the container with the id ```GID``` if it doesn't already exist, and then creates a user ```default``` in group ```GID``` with id ```UID```.
A docker compose file is provided to build and run the orgnetsim container. If you are running the container on a dev machine you can set the ```UID``` and ```GID``` as environment variables by adding the following to your ```~/.bashrc``` (or equivalent. If you do this in an existing terminal run ```source ~/.bashrc``` to re-read the ```.bashrc file```).
```
export UID=$(id -u) 
export GID=$(id -g)
```
This will set the user created in the container to the same ```GID``` and ```UID``` as your user on the host machine, allowing orgnetsim to read and write files in the directory specified in the ```volumes``` section of the ```docker.compose```.



