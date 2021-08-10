package main

import (
	"flag"
)

type Args struct {
	Verbose bool
}

func (args *Args) Parse() {
	flag.BoolVar(&args.Verbose, "v", false, "Enable verbose")
	flag.Parse()
}
