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

func NewPushWS() *restful.WebService {
	globalSender = newPushSender()

	ws := new(restful.WebService).Path("/v1/push").Consumes(restful.MIME_JSON)

	ws.Route(ws.POST("").To(globalSender.push))

	return ws
}
