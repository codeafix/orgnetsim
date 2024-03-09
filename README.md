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
This project also contains a Dockerfile that builds both the API and front-end into a container. All data that is created in the app is stored in the container path `/var/data`. By default the container will be built with an empty simulation list. To persist data outside the container make sure to mount the container path to a persistable storage path on the host machine.
```
docker run -v <host_path>:/var/data -d -p 8080:8080 orgnetsim:v0.1
```


