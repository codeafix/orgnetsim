This is the web front end for the orgnetsim utility. This website is embedded and automatically
served by the orgnetsim application, but is basically a standard REACT application.

There are several ways to run the front-end. You can use the orgnetsim docker to serve the embedded
version (default behaviour). You can override the default embedded version and specify a directory
which contains the built project, or you can run the front end using vite as below. Serving the
front end from vite will allow you to make live changes to the code.

## Available Scripts

In the `web` project directory, you can run:

### `npm start`

Runs the app in the development mode.<br>
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.<br>
You will also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.<br>
See the section about [running tests](https://vitest.dev/guide/features.html#watch-mode) for more information.

### `npm run build`

Builds the app for production to the `build` folder.<br>
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br>
Your app is ready to be deployed!

See the section about [deployment](https://vitejs.dev/guide/build) for more information.

## Learn More

You can learn more here
[Vite documentation](https://vitejs.dev/guide/)
[Vitest documentation](https://vitest.dev/guide/)

To learn React, check out the [React documentation](https://reactjs.org/).

## Using orgnetsim front end
Hopefully most of the front end features are discoverable but here are a few pointers.

### Basic features and usage
The UI opens with a list of simulations. To create a new simulation click the **Add** button and
fill in the name and description of the new simulation. Click the **Open** button to open the
simulation.

There are four sections to the simulation UI. Top left is a visualisation of the **Network**. Underneath
the network diagram on the bottom left is the **Steps...** panel which provides a visualisation of the results 
of a simulation and the table of results underneath that. Top right is the simulation **Description**, and 
underneath that on the bottom right are the **Network Options** which are used to set up a network for 
simluation.

### Network Options
This allows you to set up the initial conditions of a simulation and also modify a network as it is
imported. The default is to set all imported agents in a network to *Grey* and to have a 2 colour
simulation. By changing the **Network Options** you can also automatically add links to connect peers
within a team so that each person within each team will be connected to every other person in that
team, this is the *Link Team Peers* option. In here there are other options to automatically create
specific agents and team links to test how effective a small number of people might be in influencing
an entire network.

Also use the settings in here to set up which colours the agents can be assigned at random as they 
are created, and how many colours are available to the simulation.

### Importing a Network
To set up a simulation a network can be imported by clicking the **Import** button on the **Network** panel.
The UI will allow you to select a file to import and the settings on the **Import Network** modal control how
the network file is parsed. The system expects networks and links to be in tabular format. The minimum
required for a network is a hierarchical list of agents in tabular form. One column must contain a unique
ID, one column a parent identifier. There is an option to specify a seperate name column, but this may be 
set to the same column as the identifier.

Regular expressions may also be used to parse the information in the identifier column and extract a unique
identifier, or to match against other columns. If a regular expression is applied to another column and
it does not match the entire row will be skipped. This allows the parser to detect and skip rows that do
not contain valid data.

It is possible to load the import settings from an external json file. The json file should specify all
the settings in the following form:
```
{
    "identifier": 0,
    "name": 0,
    "parent": 2,
    "delimiter": ",",
    "regex": {
        "0":"^\\d+_(.*)$",
        "2":"^\\d+_(.*)$",
        "3":"\\S+"
    }
}
```
Note the values for the `identifier`, `name`, and `parent` settings should be without quotes so that
they are correctly interpreted as numbers and not strings. The above are the default settings on the
**Import** panel and should work to import an excel export of an organisation hierarchy from Workday.

Once a network has been imported all the nodes will appear as a tight group in the **Network** visualisation.
You can click the **Layout** button at the bottom of the **Network** panel and the agents in the network
will move to spread out of their own accord. Whilst the Network is in _Layout_ mode you can drag the nodes
around zoom in and zoom out until the Network is spread out and organised. Click the **Save Layout** button
to save this layout. Once this is done the network will remain in this layout configuration.

### Importing additional links
If you have additional data about real connections between individuals in your network you can import
them in a second step after importing the individuals in the network. To do this click the **Import**
button after importing the agents in a network. You will see that an **Import Additional Links** modal
is opened. This works in the same way as the **Import Network** modal without the option to indicate
a name column.

This modal will allow you to import tabular files that specify links between agents already on the
network by specifying the ids of the two agents to connect. Specify the column with the first Id
using the `identifier` specifier, and the column of the second Id using the `parent` specifier.
The links that are created are bi-directional so it does not matter which way round you set the
columns. Regular expressions can also be used to skip rows in the input when the expression does
not match.

This is an example settings file that will import only links that have an interaction count above
0 in the third column of the input file.
```
{
    "identifier": 0,
    "name": 0,
    "parent": 1,
    "delimiter": ",",
    "regex": {
        "2": "[1-9]\\d*"
    }
}
```

### Running Simulations
To run a simulation click the **Run** button on the **Steps...** panel. This will show the **Run Simulation**
modal. This allows you to specify the number of steps and the number of iterations per each step
to run the simulation. Each step will add a new row in the table on the **Steps...** panel. The visualisation
in the graph at the top of the **Steps...** panel shows the number of agents in each colour after all of the
conversations between agents are completed in each iteration. This is a much more granular view than the
data in the steps table.

You can run the simulation multiple times and the system will keep adding steps to the existing simulation 
for as long as you want. There is no need to keep the number of iterations the same in subsequent steps,
although the most useful number of iterations is likely to be between 10 and 100 per step.

Once there is a set of simulation results you can click the **Play** button on the **Network** panel to play
the steps through on the Network diagram so that you can see how competing ideas propagate around the 
network. This allows you to see if there are areas of the network that are more or less susceptible to 
influence as you will see these areas changing status more than other areas of the network.
