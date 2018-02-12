# sim [![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/codeafix/orgnetsim/srvr) 
This package contains all of the code for the web server that wraps the simulator

## Routes

### `GET /`
Route for the website.

### `GET /api/simulation`
Return the list of simulations on the server.

### `POST /api/simulation`
Create a new simulation.

### `PUT /api/simulation/notes`
Updates notes recorded about this list of simulations.

### `GET /api/simulation/{sim_id}`
Return a specific simulation. Contains the network options for the network if set.

### `PUT /api/simulation/{sim_id}`
Updates details about the specified simulation.

### `DELETE /api/simulation/{sim_id}`
Deletes the specified simulation

### `POST /api/simulation/{sim_id}/generate`
Generates a hierarchical network to be simulated in an existing simulation.
There should be no existing steps within the simulation otherwise this request will fail.
Returns the created first step that contains the generated network and the initial color
results for the generated network.

### `POST /api/simulation/{sim_id}/run`
Runs the simulation for a specified number of steps, each step runs a specified number of 
iterations.

### `GET /api/simulation/{sim_id}/step`
Returns the list of steps in this simulation. Currently returns the same as 
`GET /api/simulation/{sim_id}`

### `GET /api/simulation/{sim_id}/results`
Returns the concatenated set of results for all the steps in this simulation.

### `GET /api/simulation/{sim_id}/step/{step_id}`
Returns the specified step which contains the results for that step and the state of the network
at the end of that step.

### `PUT /api/simulation/{sim_id}/step/{step_id}`
Updates the results and network state in the specified step

### `DELETE /api/simulation/{sim_id}/step/{step_id}`
Deletes the specified step
