//
//  badges_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

func NewBadgesWs() *restful.WebService  {
	ws := new(restful.WebService)
	ws.Path("/badges")

	ws.Route(ws.GET("{user_id}/{sender_id}").To(badge))
	ws.Route(ws.DELETE("{user_id}/{sender_id}").To(clearBadge))

	return ws
}

func badge(request *restful.Request, response *restful.Response) {
	response.WriteHeader(http.StatusOK)
}

func clearBadge(request *restful.Request, response *restful.Response) {
	response.WriteHeader(http.StatusOK)
}
