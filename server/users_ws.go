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
	"github.com/ortuman/mercury/push"
	"github.com/ortuman/mercury/types"
	"github.com/ortuman/mercury/logger"
)

type PushSender struct {
	senders map[string]push.PushSender
}

var pushSender *PushSender

func NewUsersWS() *restful.WebService {
	pushSender = NewPushSender()

	ws := new(restful.WebService)
	ws.Path("v1/users")

	ws.Route(ws.PUT("{user_id}/{sender_id}/{token}").To(pushSender.registerToken))
	ws.Route(ws.DELETE("{user_id}/{sender_id}").To(pushSender.unregisterToken))
	ws.Route(ws.POST("{user_id}/{sender_id}").To(pushSender.push))

	return ws
}

func NewPushSender() *PushSender {
	ps := &PushSender{}
	ps.registerSenders()
	return ps
}

func (ps *PushSender) registerToken(request *restful.Request, response *restful.Response) {
	response.WriteHeader(http.StatusOK)
}

func (ps *PushSender) unregisterToken(request *restful.Request, response *restful.Response) {
	response.WriteHeader(http.StatusOK)
}

func (ps *PushSender) push(request *restful.Request, response *restful.Response) {
	var pushPayloads []types.Push
	if err := request.ReadEntity(&pushPayloads); err != nil {
		logger.Errorf("push_ws: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, payload := range pushPayloads {
		if sender, ok := ps.senders[payload.SenderID]; ok {
			sender.SendNotification(payload.UserID, payload.Notification)
		} else {
			logger.Warnf("push_ws: unrecognized sender id: %s", payload.SenderID)
		}
	}
	response.WriteHeader(http.StatusOK)
}

func (ps *PushSender) registerSenders() {
	ps.senders = make(map[string]push.PushSender)

	ps.senders[types.ApnsSenderName] = push.NewApnsSenderPool()
	ps.senders[types.GcmSenderName] = push.NewGcmSenderPool()
}
