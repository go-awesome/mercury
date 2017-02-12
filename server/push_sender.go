//
//  push_sender.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/push"
	"github.com/ortuman/mercury/logger"
)

type PushSender struct {
	senders map[string]push.PushSender
}

func NewPushSender() *PushSender {
	ps := &PushSender{
		senders: map[string]push.PushSender {
			push.ApnsSenderID : push.NewApnsSenderPool(),
			push.GcmSenderID  : push.NewGcmSenderPool(),
		},
	}
	return ps
}

func (ps *PushSender) push(request *restful.Request, response *restful.Response) {
	var pushPayloads []push.Push
	if err := request.ReadEntity(&pushPayloads); err != nil {
		logger.Errorf("push_ws: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, payload := range pushPayloads {
		for _, to := range payload.To {
			if sender, ok := ps.senders[to.SenderID]; ok {
				sender.SendNotification(&to, &payload.Notification)
			} else {
				logger.Warnf("push_ws: unrecognized sender id: %s", to.SenderID)
			}
		}
	}
	response.WriteHeader(http.StatusOK)
}
