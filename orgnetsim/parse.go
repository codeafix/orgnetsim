package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codeafix/orgnetsim/sim"
)

//OptionsFile is the format of the optional file to configure the parse and construction of a network
type OptionsFile struct {
	Network *sim.NetworkOptions `json:"network"`
	Parse   *sim.ParseOptions   `json:"parse"`
}

//Parse provides the functionality for the orgnetsim parse command utility
func Parse() {
	success, of, opt, seed := parseCommandLineOptions()
	if !success {
		return
	}

	infile := os.Args[2]
	data := readFileIntoArray(infile)

	suffix := infile[strings.LastIndex(infile, "."):len(infile)]
	of.Parse.Delimiter = ","
	if suffix == ".txt" {
		of.Parse.Delimiter = "\t"
	}

	rand.Seed(int64(seed))

	rm, err := of.Parse.ParseDelim(data)
	check(err)

	crm, err := of.Network.CloneModify(rm)
	check(err)

	n := crm.(*sim.Network)

	outfile := ""
	i := strings.LastIndex(infile, ".")
	if i > 0 {
		outfile = infile[:i] + opt + ".json"
	} else {
		outfile = infile + opt + ".json"
	}

	fo, err := os.Create(outfile)
	check(err)
	defer fo.Close()
	json := n.Serialise()
	_, err = fo.Write([]byte(json))
	check(err)
}

func parseCommandLineOptions() (success bool, of OptionsFile, opt string, seed int64) {
	of = OptionsFile{}
	of.Parse.Delimiter = ","
	of.Parse.Identifier = 0
	of.Parse.Parent = 1
	success = true
	//list of options used to change the output filename
	opt = ""

	if len(os.Args) <= 2 || os.Args[2] == "-help" {
		parsePrintUsage()
		return false, of, opt, seed
	}

	//List of unrecognised command switches
	uc := []string{}

	skipnext := false
	for i, arg := range os.Args[3:len(os.Args)] {
		if skipnext {
			skipnext = false
			continue
		}
		switch arg {
		case "-awm":
			of.Network.AgentsWithMemory = true
			opt = opt + arg
		case "-ltp":
			of.Network.LinkTeamPeers = true
			opt = opt + arg
		case "-ic":
			of.Network.InitColors = []sim.Color{sim.Grey, sim.Red}
			opt = opt + arg
		case "-be":
			opt = opt + arg
			if len(os.Args) < i+5 {
				fmt.Printf("<beListFile> missing after -be option \n\n")
				success = false
				break
			}
			beFile := os.Args[i+4]
			skipnext = true
			of.Network.EvangelistList = readFileIntoArray(beFile)
		case "-lt":
			opt = opt + arg
			if len(os.Args) < i+5 {
				fmt.Printf("<ltListFile> missing after -lt option \n\n")
				success = false
				break
			}
			ltFile := os.Args[i+4]
			skipnext = true
			of.Network.LinkedTeamList = readFileIntoArray(ltFile)
		case "-mc":
			opt = opt + arg
			if len(os.Args) < i+5 {
				fmt.Printf("<maxColors> missing after -mc option \n\n")
				success = false
				break
			}
			skipnext = true
			var err error
			of.Network.MaxColors, err = strconv.Atoi(os.Args[i+4])
			if err != nil {
				fmt.Printf("Invalid integer '%s' for -mc option\n", os.Args[i+4])
				success = false
				break
			}
			opt = opt + os.Args[i+4]
		case "-opt":
			opt = opt + arg
			if len(os.Args) < i+5 {
				fmt.Printf("<optionsFile> missing after -opt option \n\n")
				success = false
				break
			}
			skipnext = true
			file := os.Args[i+4]
			optFile := strings.Join(readFileIntoArray(file), "")
			err := json.Unmarshal([]byte(optFile), &of)
			if err != nil {
				fmt.Printf("Error in <optionsFile>: %s \n\n", err.Error())
				success = false
				break
			}
			s := strings.Split(file, "\\")
			l := strings.Split(s[len(s)-1], "/")
			opt = opt + "-" + l[len(l)-1]
		case "-seed":
			opt = opt + arg
			if len(os.Args) < i+5 {
				fmt.Printf("<seed> missing after -seed option \n\n")
				success = false
				break
			}
			skipnext = true
			var err error
			seed, err = strconv.ParseInt(os.Args[i+4], 10, 64)
			if err != nil {
				fmt.Printf("Invalid integer '%s' for -seed option\n", os.Args[i+4])
				success = false
				break
			}
			opt = opt + os.Args[i+4]
		default:
			uc = append(uc, arg)
		}
	}
	if of.Network.MaxColors == 0 {
		//MaxColors has not been set so default to 4
		of.Network.MaxColors = 4
		opt = opt + "-mc4"
	}
	if seed == 0 {
		//seed has not been set so default to time.Now
		seed = time.Now().UnixNano()
	}
	if len(uc) > 0 {
		fmt.Printf("Unrecognised options on command line: %s\n\n", strings.Join(uc, " "))
		success = false
	}
	return success, of, opt, seed
}

func readFileIntoArray(filename string) []string {
	r := []string{}
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		r = append(r, s.Text())
	}
	return r
}

func parsePrintUsage() {
	fmt.Println("Converts csv or tab delimited text into an orgnetsim network")
	fmt.Println("NOTE: This utility only supports utf-8 encoded files so if you saved from Excel you will")
	fmt.Println("      have to convert the file format.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("      orgnetsim parse <orglist> [-opt <optionsFile>] [-awm] [-ltp] [-ic] [-be <beListFile>] [-lt <ltListFile>] [-seed <seed>] [-mc <maxColors>] [-seed <seed>]")
	fmt.Println("      orgnetsim parse -help")
	fmt.Println()
	fmt.Println("<orglist>")
	fmt.Println("      is a file that contains the list of individuals in an organisation.")
	fmt.Println("      The default delimiter is a \",\".")
	fmt.Println("      The first line in the file is assumed to be column headers and is skipped.")
	fmt.Println("      Each line is understood as a single individual with the first column being a unique")
	fmt.Println("      identifier and the second column containing the unique identifier of the individuals")
	fmt.Println("      direct parent. All other information is ignored.")
	fmt.Println("-awm")
	fmt.Println("      Use agents with memory. Default is off.")
	fmt.Println("-ltp")
	fmt.Println("      Add links that connect each member of a team with every other member of that team.")
	fmt.Println("      Default is off.")
	fmt.Println("-ic")
	fmt.Println("      Randomly assign a Color to the agents in the network either Grey or Red. Default all")
	fmt.Println("      individuals are Grey.")
	fmt.Println("-be <beListFile>")
	fmt.Println("      Alter all the individuals listed in the beListFile, set their Color to Blue, and")
	fmt.Println("      their susceptibility scores artificially high so they remain evangelists for a")
	fmt.Println("      specific idea.")
	fmt.Println("-lt <ltListFile>")
	fmt.Println("      Connect all the individuals listed in the ltListFile to each other as a single team.")
	fmt.Println("-seed <seed>")
	fmt.Println("      The seed for generating random properties of agents. Can be set to any integer value.")
	fmt.Println("      The default is time.Now() in nanoseconds.")
	fmt.Println("-mc <maxColors>")
	fmt.Println("      The maximum number of colors permitted on the network in the simulation. Default")
	fmt.Println("      is 4.")
	fmt.Println("-opt <optionsFile>")
	fmt.Println("      A file containing the ParseOptions and NetworkOptions to apply when parsing the <orglist>.")
	fmt.Println("      The file should be in the following format. Any options not present are given their")
	fmt.Println("      default value. In the Regex option you can supply regular expressions to filter data rows")
	fmt.Println("      in the <orglist>. If the regular expressions do not match the content of the column which")
	fmt.Println("      they are applied to (specified by the integer column index) then the row will be ignored.")
	fmt.Println("      {")
	fmt.Println("        \"network\":{")
	fmt.Println("          \"linkTeamPeers\": false,")
	fmt.Println("          \"linkedTeamList\": [\"id_1\",\"id_2\"],")
	fmt.Println("          \"evangelistList\": [\"id_1\",\"id_2\"],")
	fmt.Println("          \"loneEvangelist\": [\"id_1\",\"id_2\"],")
	fmt.Println("          \"initColors\": [0,2],")
	fmt.Println("          \"maxColors\": 4,")
	fmt.Println("          \"agentsWithMemory\": false")
	fmt.Println("        },")
	fmt.Println("        \"parse\":{")
	fmt.Println("          \"identifier\": 0,")
	fmt.Println("          \"parent\": 1,")
	fmt.Println("          \"delimiter\": \",\",")
	fmt.Println("          \"regex\": {\"2\":\"\\\\S+\"}")
	fmt.Println("        }")
	fmt.Println("      }")
	fmt.Println("-help")
	fmt.Println("      Prints this message.")
}
