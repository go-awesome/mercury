//
//  server/push_sender.go
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

type pushSender struct {
	senders map[string]*push.SenderHub
}

func newPushSender() *pushSender {
	ps := &pushSender{
		senders: map[string]*push.SenderHub{
			push.ApnsSenderID    : push.NewApnsSenderHub(),
			push.GcmSenderID     : push.NewGcmSenderHub(),
			push.SafariSenderID  : push.NewSafariSenderHub(),
			push.ChromeSenderID  : push.NewChromeSenderHub(),
			push.FirefoxSenderID : push.NewFirefoxSenderHub(),
		},
	}
	return ps
}

func (ps *pushSender) push(request *restful.Request, response *restful.Response) {
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

func (ps *pushSender) Stats() map[string]push.PushStats {
	stats := map[string]push.PushStats{}

	for senderID, sender := range ps.senders {
		stats[senderID] = sender.Stats()
	}
	return stats
}
