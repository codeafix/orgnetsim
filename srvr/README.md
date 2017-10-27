# sim [![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/codeafix/orgnetsim/srvr) 
This package contains all of the code for the web server that wraps the simulator

## Routes

### `GET /`
Route for the website

### `GET /api/simulation`
Return the list of simulations on the server

### `POST /api/simulation`
Create a new simulation

### `GET /api/simulation/{sim_id}`
Return a specific simulation

### `PUT /api/simulation/{sim_id}`
Updates details about the specified simulation

### `GET /api/simulation/{sim_id}/network`
Returns the current state of the network being simulated.

### `PUT /api/simulation/{sim_id}/network`
Updates the current state of the network being simulated

### `POST /api/simulation/{sim_id}/generate`
Generates a hierarchical network to be simulated

### `POST /api/simulation/{sim_id}/run`
Runs the simulation for a specified number of steps, each step runs a specified number of iterations

### `GET /api/simulation/{sim_id}/step`
Returns the list of steps in this simulation. There will always be a step 0 which
contains the initial conditions of the simulation.

### `GET /api/simulation/{sim_id}/step/{step_id}`
Returns the results for the given step

### `GET /api/simulation/{sim_id}/step/{step_id}/network`
Returns the state of the network at the end of the given step

### `GET /api/simulation/{sim_id}/results`
Returns the concatenated set of results for all the steps in this simulation
