//
//  ping_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"github.com/emicklei/go-restful"
)

var globalSender *pushSender

func init() {
	globalSender = newPushSender()
}

func NewPushWS() *restful.WebService {
	ws := new(restful.WebService).Path("/v1/push").Consumes(restful.MIME_JSON)

	ws.Route(ws.POST("").To(globalSender.push))
	ws.Route(ws.GET("/stats").To(globalSender.stats))

	return ws
}
