//
//  main.go
//  mercuryx
//
//  Copyright (c) 2016 Miguel Ángel Ortuño. All rights reserved.
//

package main

import (
	"os"
	"github.com/ortuman/mercury/server"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
)

const configFile = "/etc/mercury/mercury.conf"

func main() {
	// load configuration
	if _, err := os.Stat(configFile); err == nil {
		config.Load(configFile)
	} else {
		logger.Warnf("main: couldn't load config file '%s': %v", configFile, err)
	}

	srv := server.NewServer()
	srv.Run()
}
