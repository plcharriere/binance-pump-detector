package main

import (
	"flag"
)

type Args struct {
	Verbose        bool
	ConfigFilePath string
}

func (args *Args) Parse() {
	flag.BoolVar(&args.Verbose, "v", false, "Enable verbose")
	flag.StringVar(&args.ConfigFilePath, "f", "config.ini", "Configuration file path")
	flag.Parse()
}
