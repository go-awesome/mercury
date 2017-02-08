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
	logger.Infof("server: accepting commands at %s:%s [%s]", config.Server.ListAddr, cmdPort, runtime.Version())

	/*
	[PUT]    /users/34/gcm
	[DELETE] /users/34/gcm

	[POST]   /users/34/apns

	[GET]    /users/apns/gone
	[DELETE] /users/apns/gone

	[GET]    /badges/34/gcm
	[DELETE] /badges/34/gcm
	*/

	// configure services
	restful.Add(NewPingWS())   /* /ping */
	restful.Add(NewUsersWS())  /* /users */
	restful.Add(NewBadgesWs()) /* /badges */

	addr := config.Server.ListAddr + ":" + cmdPort

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Errorf("server: ListenAndServe: %v", err)
	}
}
