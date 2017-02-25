//
//  push/safari.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"github.com/ortuman/mercury/cert"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
	"github.com/timehop/apns"
	"time"
)

func NewSafariSenderHub() *SenderHub {
	s := NewSenderHub(SafariSenderID, NewSafariSender, config.Safari.MaxConn)
	logger.Infof("safari: initialized safari sender (max conn: %d)", s.senderCount)
	return s
}

// MARK: SafariPushSender

type SafariPushSender struct {
	client *apns.Client
}

func NewSafariSender() (PushSender, error) {
	s := &SafariPushSender{}
	certificate, key, err := cert.LoadP12(config.Safari.CertFile, "")
	if err != nil {
		return nil, err
	}
	tlsCert := cert.TLSCertificate(certificate, key)

	client := apns.NewClientWithCert(apns.ProductionGateway, tlsCert)
	s.client = &client

	return s, nil
}

func (ss *SafariPushSender) SendNotification(to *To, notification *Notification) (int, time.Duration) {

	p := apns.NewPayload()
	p.APS.Alert.Title = notification.Title
	p.APS.Alert.Body = notification.Body
	p.APS.Alert.Action = notification.Action
	p.APS.URLArgs = notification.URLArgs

	expiration := time.Now().Add(time.Hour * 24 * 7)

	m := apns.NewNotification()
	m.ID = notification.ID
	m.DeviceToken = to.DeviceToken
	m.Priority = apns.PriorityImmediate
	m.Expiration = &expiration
	m.Payload = p

	start := time.Now().UnixNano()
	err := ss.client.Send(m)
	end := time.Now().UnixNano()
	reqElapsed := time.Duration((end - start)) / time.Millisecond

	if err == nil {
		logger.Debugf("safari: [%s] notification delivered: %s", to.UserID, notification.ID)
		return StatusDelivered, reqElapsed

	} else {
		logger.Errorf("safari: [%s] notification COULDN'T be delivered: %s: %v", to.UserID, notification.ID, err)
		return StatusFailed, 0
	}
}
