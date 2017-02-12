//
//  server.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"strconv"
	"runtime"
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/storage"
)

type Server struct {
}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Run() {

	// configure logger
	logger.SetLogFilePath(config.Logger.Logfile)
	logger.SetLogLevel(config.Logger.Level)

	// initialize storage
	storage.Instance()

	// listen commands
	cmdPort := strconv.Itoa(config.Server.Port)

	logger.Infof("server: %s %s", config.ServiceName, config.ServiceVersion)
	logger.Infof("server: accepting commands at %s:%s [%s]", config.Server.BindAddress, cmdPort, runtime.Version())

	// configure services
	restful.Add(NewPingWS())   /* /ping */
	restful.Add(NewPushWS())   /* /push */
	restful.Add(NewBadgesWS()) /* /badges */

	addr := config.Server.BindAddress + ":" + cmdPort

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Errorf("server: ListenAndServe: %v", err)
	}
}
