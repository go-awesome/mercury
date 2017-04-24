//
//  server/server.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/storage"
	"net/http"
	"strconv"
	"os"
	"path/filepath"
	"fmt"
	"strings"
	"net"
)

type Server struct {
}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Run() error {

	// configure logger
	logger.SetLogFilePath(config.Logger.Logfile)
	logger.SetLogLevel(config.Logger.Level)

	// create PID file
	if err := s.createPIDFile(); err != nil {
		logger.Warnf("server: pid: %v", err)
	}

	// initialize storage
	storage.Instance()

	// start server...
	logger.Infof("server: %s %s", config.ServiceName, config.ServiceVersion)

	if !strings.HasPrefix(config.Server.ListenAddr, "unix:") {
		return s.listenOnTCPSocket()
	} else {
		return s.listenOnUnixDomainSocket()
	}
}

func (s *Server) createPIDFile() error {
	if !config.PID.Enabled {
		return nil
	}

	pidPath := config.PID.File
	_, err := os.Stat(pidPath)
	if os.IsNotExist(err) || config.PID.Override {
		currentPid := os.Getpid()
		if err := os.MkdirAll(filepath.Dir(pidPath), os.ModePerm); err != nil {
			return err
		}

		file, err := os.Create(pidPath)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(strconv.FormatInt(int64(currentPid), 10)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s already exists", pidPath)
	}
	return nil
}

func (s *Server) configureServices() {
	restful.Add(NewPingWS()) // /ping

	restful.Add(NewPushWS())     // /v1/push
	restful.Add(NewBadgesWS())   // /v1/badges
	restful.Add(NewStatsWS())    // /v1/stats
	restful.Add(NewSendFormWS()) // /v1/send_form
}

func (s *Server) listenOnTCPSocket() error {
	logger.Infof("server: accepting commands at: %s", config.Server.ListenAddr)
	s.configureServices()
	return http.ListenAndServe(config.Server.ListenAddr, nil)
}

func (s *Server) listenOnUnixDomainSocket() error {
	unixServicePath := strings.TrimPrefix(config.Server.ListenAddr, "unix:")

	// unlink previous sock file if existing
	os.Remove(unixServicePath)

	if err := os.MkdirAll(filepath.Dir(unixServicePath), os.ModePerm); err != nil {
		return err
	}

	addr, err := net.ResolveUnixAddr("unix", unixServicePath)
	if err != nil {
		return err
	}

	l, err := net.ListenUnix("unix", addr)
	if err != nil {
		return err
	}

	logger.Infof("server: accepting commands at unix domain socket: %s", unixServicePath)
	s.configureServices()
	return http.Serve(l, nil)
}
