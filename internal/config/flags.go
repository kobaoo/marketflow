package config

import (
	"flag"
	"fmt"
)

var Port int

func ParseFlags() {
	flag.IntVar(&Port, "port", 0, "Port to serve on")
	flag.Usage = printHelp

	flag.Parse()
}

func printHelp() {
	fmt.Print(`Usage:
  marketflow [--port <N>]
  marketflow --help

Options:
  --port N     Port number
  `)
}
