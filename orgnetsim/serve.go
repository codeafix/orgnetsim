package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codeafix/orgnetsim/srvr"
)

//ServeOptions holds settings specified on the command line for the serve command
type ServeOptions struct {
	Port string
}

//Serve provides the functionality for the orgnetsim serve command utility
func Serve() {
	success, so := serveCommandLineOptions()
	if !success {
		return
	}

	srvr.ListenAndServe(os.Args[2], so.Port)
}

func serveCommandLineOptions() (success bool, so ServeOptions) {
	so = ServeOptions{}
	so.Port = "8080"
	success = true

	if len(os.Args) <= 2 || os.Args[2] == "-help" {
		servePrintUsage()
		return false, so
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
		case "-p":
			if len(os.Args) < i+5 {
				fmt.Printf("<port> missing after -p option \n\n")
				success = false
				break
			}
			so.Port = os.Args[i+4]
			skipnext = true
		default:
			uc = append(uc, arg)
		}
	}
	if len(uc) > 0 {
		fmt.Printf("Unrecognised options on command line: %s\n\n", strings.Join(uc, " "))
		success = false
	}
	return success, so
}

func servePrintUsage() {
	fmt.Println("Starts an orgnetsim server that persists simulations in the folder specified by <rootpath>.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("      orgnetsim serve <rootpath> [-p <port>]")
	fmt.Println("      orgnetsim serve -help")
	fmt.Println()
	fmt.Println("<rootpath>")
	fmt.Println("      is a folder where the server will store all resources that are created and updated by the")
	fmt.Println("      orgnetsim routes")
	fmt.Println("-p")
	fmt.Println("      Specifies the port that the server will listen on. The default is 8080.")
	fmt.Println("-help")
	fmt.Println("      Prints this message.")
}
