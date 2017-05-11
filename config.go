package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	configfile = flag.String("conf", "conftest.toml", "path for configuration toml file")
)

type Config struct {
	Jsonpaths string
	Logfile   string
	Wait      int
}

// Reads info from config file
func readConfig() Config {
	_, err := os.Stat(*configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(*configfile, &config); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return config
}
