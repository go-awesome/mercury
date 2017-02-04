//
//  push_ws.go
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

func NewPushWS() *restful.WebService {
	pushSender = NewPushSender()

	s := new(restful.WebService)
	s.Path("/push")
	s.Consumes(restful.MIME_JSON)
	s.Route(s.POST("").To(pushSender.push))

	return s
}

func NewPushSender() *PushSender {
	ps := &PushSender{}
	ps.registerSenders()
	return ps
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
			sender.SendNotification(payload.UserID, payload.Notification, payload.Auth)
		} else {
			logger.Warnf("push_ws: unrecognized sender id: %s", payload.SenderID)
		}
	}
	response.WriteHeader(http.StatusOK)
}

func (ps *PushSender) registerSenders() {
	ps.senders = make(map[string]push.PushSender)

	ps.senders[push.ApnsSenderID] = push.NewApnsSenderPool()
	ps.senders[push.GcmSenderID] = push.NewGcmSenderPool()
}
