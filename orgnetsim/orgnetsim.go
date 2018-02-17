package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 || os.Args[1] == "-help" {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "parse":
		Parse()
	case "serve":
		Serve()
	default:
		fmt.Printf("Unrecognised command line: %s\n\n", os.Args[1])
		printUsage()
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printUsage() {
	fmt.Println("orgnetsim command line utility")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("      orgnetsim <command> [options]*")
	fmt.Println("      orgnetsim -help")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    parse <orglist> [-help] [-awm] [-ltp] [-ic] [-be <beListFile>] [-lt <ltListFile>] [-mc <maxColors>]")
	fmt.Println("        Reads in a csv or tsv and converts into an orgnetsim network saved in json format.")
	fmt.Println("    serve <rootpath> [-help] [-p <port>]")
	fmt.Println("        Starts an orgnetsim server that persists simulations in the folder specified by <rootpath>.")
	fmt.Println("-help")
	fmt.Println("    Prints this message.")
}
