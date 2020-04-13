package main

import (
	"flag"
	"fmt"
	"os"
)

func printDefaults() {
	fmt.Printf("Enter a querystring or available commands: 'init', 'help' \n")
	os.Exit(1)
}

// Cli does command line interface
func Cli() {

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	golinksHostname := initCmd.String("hostname", "", "hostname")
	golinksPort := initCmd.String("port", "80", "port")
	golinksProtocol := initCmd.String("protocol", "http", "protocol")

	if len(os.Args) < 2 {
		printDefaults()
	}

	switch os.Args[1] {

	case "init":
		initCmd.Parse(os.Args[2:])
		if initCmd.Parsed() {
			// Required Flags
			if *golinksHostname == "" {
				initCmd.PrintDefaults()
				os.Exit(1)
			}
			if *golinksPort == "" {
				initCmd.PrintDefaults()
				os.Exit(1)
			}
			if *golinksProtocol == "" {
				initCmd.PrintDefaults()
				os.Exit(1)
			}
		}

		initialize(*golinksHostname, *golinksPort, *golinksProtocol)

	case "--help":
		printDefaults()

	default:
		if os.Args[1] == "" {
			printDefaults()
		}
		queryBrowse(os.Args[1])
	}

}
