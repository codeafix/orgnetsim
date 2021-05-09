# orgnetsim command line utility
Usage:
```
      orgnetsim <command> [options]*
      orgnetsim -help
```
Commands:
```
    parse <orglist> [-help] [-opt <optionsFile>] [-awm] [-ltp] [-ic] [-be <beListFile>] [-lt <ltListFile>] [-mc <maxColors>]
```

Reads in a csv or tsv and converts into an orgnetsim network saved in json format.
```
    serve <rootpath> <webdir> [-p <port>]
```

Starts an orgnetsim server that persists simulations in the folder specified by <rootpath>.
Serves static web files from the folder specified by <webdir>.

`-help`
Prints this message.

Converts csv or tab delimited text into an orgnetsim network
NOTE: This utility only supports utf-8 encoded files so if you saved from Excel you will
      have to convert the file format.

## orgnetsim parse
Usage:
```
      orgnetsim parse <orglist> [-opt <optionsFile>] [-awm] [-ltp] [-ic] [-be <beListFile>] [-lt <ltListFile>] [-seed <seed>] [-mc <maxColors>] [-seed <seed>]
      orgnetsim parse -help
```

`<orglist>`
is a file that contains the list of individuals in an organisation.
Comma separated files are supported and must have *.csv suffix.
Tab separated files are supported and must have *.txt suffix.
The first line in the file is assumed to be column headers and is skipped.
Each line is understood as a single individual with the first column being a unique
identifier and the second column containing the unique identifier of the individuals
direct parent. All other information is ignored.

`-awm`
Use agents with memory. Default is off.

`-ltp`
Add links that connect each member of a team with every other member of that team.
Default is off.

`-ic`
Randomly assign a Color to the agents in the network either Grey or Red. Default all
individuals are Grey.

`-be <beListFile>`
Alter all the individuals listed in the beListFile, set their Color to Blue, and
their susceptibility scores artificially high so they remain evangelists for a
specific idea.

`-lt <ltListFile>`
Connect all the individuals listed in the ltListFile to each other as a single team.

`-seed <seed>`
The seed for generating random properties of agents. Can be set to any integer value.
The default is time.Now() in nanoseconds.

`-mc <maxColors>`
The maximum number of colors permitted on the network in the simulation. Default
is 4.

`-opt <optionsFile>`
A file containing the ParseOptions and NetworkOptions to apply when parsing the `<orglist>`.
The file should be in the following format. Any options not present are given their
default value. In the Regex option you can supply regular expressions to filter data in
the rows in the `<orglist>`. If the regular expressions do not match the content of the column
which they are applied to (specified by the integer column index) then row will be ignored.
```
      {
        "network":{
          "linkTeamPeers": false,
          "linkedTeamList": ["id_1","id_2"],
          "evangelistList": ["id_1","id_2"],
          "loneEvangelist": ["id_1","id_2"],
          "initColors": [0,3],
          "maxColors": 4,
          "agentsWithMemory": false
        },
        "parse":{
          "identifier": 0,
          "parent": 1,
          "delimiter": ",",
          "regex": {"2":"\\S+"}
        }
      }
```

`-help`
Prints this message.

## orgnetsim serve
Usage:
```
      orgnetsim serve <rootpath> <webdir> [-p <port>]
      orgnetsim serve -help
```

`<rootpath>`
is a folder where the server will store all resources that are created and updated by the
orgnetsim routes

`<webdir>`
is a folder where the static files for a web front end are served from.

`-p`
Specifies the port that the server will listen on. The default is 8080.

`-help`
Prints this message.