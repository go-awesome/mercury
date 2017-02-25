//
//  mercury.go
//  mercury
//
//  Copyright (c) 2016 Miguel Ángel Ortuño. All rights reserved.
//

package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/ortuman/mercury/server"
	"github.com/ortuman/mercury/config"
)

func main() {

	var (
		configFile string
		help bool
	)

	flag.StringVar(&configFile, "config", "/etc/mercury/mercury.conf", "configuration path file")
	flag.BoolVar(&help, "help", false, "show application usage")
	flag.Parse()

	if !help {
		// load configuration
		if _, err := os.Stat(configFile); err == nil {
			config.Load(configFile)
		} else {
			fmt.Fprintf(os.Stderr, "mercury: couldn't load config file '%s': %v\n", configFile, err)
			os.Exit(-1)
		}

		srv := server.NewServer()
		srv.Run()
	} else {
		flag.Usage()
	}
}
