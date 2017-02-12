//
//  mercury.go
//  mercuryx
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

	var configFile string
	flag.StringVar(&configFile, "config-file", "/etc/mercury/mercury.conf", "configuration path file")
	flag.Parse()

	// load configuration
	if _, err := os.Stat(configFile); err == nil {
		config.Load(configFile)
	} else {
		fmt.Fprintf(os.Stderr, "mercury: couldn't load config file '%s': %v", configFile, err)
		os.Exit(-1)
	}

	srv := server.NewServer()
	srv.Run()
}
