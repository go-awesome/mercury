//
//  server/ping_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/logger"
	"net/http"
)

func NewPingWS() *restful.WebService {
	ws := new(restful.WebService).Path("/ping")
	ws.Route(ws.GET("").To(ping))
	return ws
}

// Checks if the server is alive. This is useful for monitoring tools, load-balancers and automated scripts.
func ping(_ *restful.Request, response *restful.Response) {
	logger.Infof("ping_ws: pong...")
	response.WriteHeader(http.StatusOK)
}
