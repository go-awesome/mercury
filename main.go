//
//  main.go
//  mercuryx
//
//  Copyright (c) 2016 Miguel Ángel Ortuño. All rights reserved.
//

package main

import (
	"os"
	"github.com/Hooks-Alerts/mercuryx/server"
	"github.com/Hooks-Alerts/mercuryx/config"
	"github.com/Hooks-Alerts/mercuryx/logger"
)

const configFile = "/etc/mercuryx/mercuryx.conf"

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
