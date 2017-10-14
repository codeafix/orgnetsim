# orgnetsim [![Build Status](https://travis-ci.org/codeafix/orgnetsim.svg?branch=master)](https://travis-ci.org/codeafix/orgnetsim) [![Coverage Status](http://codecov.io/github/codeafix/orgnetsim/coverage.svg?branch=master)](http://codecov.io/github/codeafix/orgnetsim?branch=master) [![MIT](https://img.shields.io/npm/l/express.svg)](https://github.com/codeafix/orgnetsim/blob/master/LICENSE)
A simulator for Organisational Networks

The simulator is created from a Network of Agents. The Network itself can be any arbitrary graph and contains a collection of Agents and a collection of links between those Agents. The simulator uses Colors to represent competing ideas on the Network. The default Color for an Agent is Grey. During a simulation Agents interact and decide whether or not to update their Color.

## Packages

[sim](sim/README.md) The organisation network simulator

## TODOs
- [ ] Integrate into a web service
- [ ] Create a network visualiser using D3