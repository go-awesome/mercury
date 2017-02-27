//
//  server/badges_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
    "fmt"
    "net/http"
    "github.com/emicklei/go-restful"
    "github.com/ortuman/mercury/storage"
    "github.com/ortuman/mercury/logger"
    "github.com/ortuman/mercury/push"
)

func NewBadgesWS() *restful.WebService {
    ws := new(restful.WebService).Path("/v1/badges")

    ws.Route(ws.GET("/{sender-id}/{token}").To(badge))
    ws.Route(ws.DELETE("/{sender-id}/{token}").To(clearBadge))

    return ws
}

func badge(request *restful.Request, response *restful.Response) {
    senderID := request.PathParameter("sender-id")
    token := request.PathParameter("token")

    if !push.IsValidSenderID(senderID) {
        logger.Warnf("badges_ws: unrecognized sender id: %s", senderID)
        response.WriteHeader(http.StatusBadRequest)
        return
    }

    badge, err := storage.Instance().GetBadge(senderID, token)
    if err != nil {
        logger.Errorf("badges_ws: %v", err)
        response.WriteHeader(http.StatusInternalServerError)
        return
    }
    response.Write([]byte(fmt.Sprintf("%d", badge)))
}

func clearBadge(request *restful.Request, response *restful.Response) {
    senderID := request.PathParameter("sender-id")
    token := request.PathParameter("token")

    if !push.IsValidSenderID(senderID) {
        logger.Warnf("badges_ws: unrecognized sender id: %s", senderID)
        response.WriteHeader(http.StatusBadRequest)
        return
    }

    err := storage.Instance().ClearBadge(senderID, token)
    if err != nil {
        logger.Errorf("badges_ws: %v", err)
        response.WriteHeader(http.StatusInternalServerError)
        return
    }
    response.WriteHeader(http.StatusOK)
}
