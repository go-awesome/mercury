//
//  ping_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/logger"
)

func NewPingWS() *restful.WebService {
	s := new(restful.WebService)
	s.Path("/ping")
	s.Route(s.GET("/*").To(ping))
	return s
}

// Checks if the server is alive. This is useful for monitoring tools, load-balancers and automated scripts.
func ping(_ *restful.Request, response *restful.Response) {
	logger.Infof("ping_ws: pong...")
	response.WriteHeader(http.StatusOK)
}
