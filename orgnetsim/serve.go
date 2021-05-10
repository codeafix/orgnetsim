package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/codeafix/orgnetsim/srvr"
)

//ServeOptions holds settings specified on the command line for the serve command
type ServeOptions struct {
	Port string
}

//Serve provides the functionality for the orgnetsim serve command utility
func Serve(webfs fs.FS) {
	success, staticDir, so := serveCommandLineOptions()
	if !success {
		return
	}

	srvr.ListenAndServe(os.Args[2], staticDir, webfs, so.Port)
}

func serveCommandLineOptions() (success bool, staticDir string, so ServeOptions) {
	so = ServeOptions{}
	so.Port = "8080"
	staticDir = ""
	success = true

	if len(os.Args) < 3 || os.Args[2] == "-help" {
		servePrintUsage()
		return false, staticDir, so
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
		case "-s":
			if len(os.Args) < i+5 {
				fmt.Printf("<staticDir> missing after -s option \n\n")
				success = false
				break
			}
			staticDir = os.Args[i+4]
			skipnext = true
		default:
			uc = append(uc, arg)
		}
	}
	if len(uc) > 0 {
		fmt.Printf("Unrecognised options on command line: %s\n\n", strings.Join(uc, " "))
		success = false
	}
	return success, staticDir, so
}

func servePrintUsage() {
	fmt.Println("Starts an orgnetsim server that persists simulations in the folder specified by <rootpath>")
	fmt.Println("and serves the website in the folder specified by <webpath>.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("      orgnetsim serve <rootpath> [-s <webpath>] [-p <port>]")
	fmt.Println("      orgnetsim serve -help")
	fmt.Println()
	fmt.Println("<rootpath>")
	fmt.Println("      is a folder where the server will store all resources that are created and updated by the")
	fmt.Println("      orgnetsim routes.")
	fmt.Println("-s <webpath>")
	fmt.Println("      is a folder containing the static website that is served by the server. By default the")
	fmt.Println("      server will use a copy of the website embedded in the executable.")
	fmt.Println("-p <port>")
	fmt.Println("      Specifies the port that the server will listen on. The default is 8080.")
	fmt.Println("-help")
	fmt.Println("      Prints this message.")
}
