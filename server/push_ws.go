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
)

func NewPushWS() *restful.WebService {
	ws := new(restful.WebService).Path("/v1/push").Consumes(restful.MIME_JSON)

	ws.Route(ws.POST("").To(sendPush))

	return ws
}

func sendPush(request *restful.Request, response *restful.Response) {
	response.WriteHeader(http.StatusOK)
}
