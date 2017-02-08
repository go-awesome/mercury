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
	"github.com/ortuman/mercury/types"
	"github.com/ortuman/mercury/push"
	"github.com/ortuman/mercury/logger"
)

type PushSender struct {
	senders map[string]push.PushSender
}

func NewPushSender() *PushSender {
	ps := &PushSender{
		senders: map[string]push.PushSender {
			types.ApnsSenderID : push.NewApnsSenderPool(),
			types.GcmSenderID  : push.NewGcmSenderPool(),
		},
	}
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
		for _, userID := range payload.UserIDs {
			if sender, ok := ps.senders[payload.SenderID]; ok {
				sender.SendNotification(userID, &payload.Notification)
			} else {
				logger.Warnf("push_ws: unrecognized sender id: %s", payload.SenderID)
			}
		}
	}
	response.WriteHeader(http.StatusOK)
}
