//
//  users_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/storage"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/types"
)

var pushSender *PushSender

func NewUsersWS() *restful.WebService {
	pushSender = NewPushSender()

	ws := new(restful.WebService)
	ws.Path("v1/users")

	ws.Route(ws.PUT("{user_id}/{sender_id}/{token}").To(registerToken))
	ws.Route(ws.DELETE("{user_id}/{sender_id}").To(unregisterToken))

	ws.Route(ws.POST("{user_id}/{sender_id}").To(pushSender.push))

	return ws
}

func registerToken(request *restful.Request, response *restful.Response) {
	userID := request.PathParameter("user_id")
	senderID := request.PathParameter("sender_id")
	token := request.PathParameter("token")

	if types.IsValidSenderID(senderID) {
		logger.Errorf("users_ws: invalid sender id: %s", senderID)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	si := &storage.SenderInfo{
		UserID: 	userID,
		SenderID: 	senderID,
		Token: 		token,
	}
	if err := storage.Instance().InsertSenderInfo(si); err != nil {
		logger.Errorf("users_ws: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func unregisterToken(request *restful.Request, response *restful.Response) {
	userID := request.PathParameter("user_id")
	senderID := request.PathParameter("sender_id")

	if types.IsValidSenderID(senderID) {
		logger.Errorf("users_ws: invalid sender id: %s", senderID)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := storage.Instance().DeleteSenderInfo(userID, senderID); err != nil {
		logger.Errorf("users_ws: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}
