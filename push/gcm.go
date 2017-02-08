//
//  gcm.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"net/http"
	"golang.org/x/net/http2"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/types"
)

const gcmSendEndpoint = "https://android.googleapis.com/gcm/send"

const gcmNotRegisteredError = "NotRegistered"

func NewGcmSenderPool() *SenderPool {
	s := NewSenderPool("gcm", NewGcmPushSender, config.Gcm.PoolSize)
	logger.Infof("gcm: initialized gcm sender (pool size: %d)", s.senderCount)
	return s
}

// MARK: GcmPushSender

type GcmPushSender struct {
	client *http.Client
}

func NewGcmPushSender() (PushSender, error) {
	s := &GcmPushSender{}
	transport := &http.Transport{}
	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
	}
	s.client = &http.Client{Transport: transport}
	return s, nil
}

func (s *GcmPushSender) SendNotification(userID string, notification *types.Notification) {

	/*
	// prepare GCM request & response entities
	alert := make(map[string]string)
	alert["title"] = gcmNotification.Title
	alert["body"] = gcmNotification.Body

	data := make(map[string]interface{})
	data["alert"] = alert
	data["notification"] = notification
	data["id"] = gcmNotification.Identifier
	data["title"] = gcmNotification.Title
	data["action"] = "open_notification"

	gcmData := make(map[string]interface{})
	gcmData["priority"] = "high"
	gcmData["data"] = data

	regIDs := []string{gcmAuth.RegistrationID}
	gcmReq := &GcmRequest{RegistrationIDs: regIDs, Data: gcmData}
	gcmReq.TimeToLive = 86400

	gcmResp := &GcmResponse{}

	// perform request
	r := NewRequest(s.client)
	defer r.Close()

	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("key=%s", gcmAuth.ApiKey)

	statusCode, err := r.URL(gcmSendEndpoint).Headers(headers).POST().Do(&gcmResp)
	if err != nil {
		logger.Errorf("gcm: %v", err)
		return
	}

	if statusCode == http.StatusOK {
		for _, result := range gcmResp.Results {
			if len(result.Error) == 0 {
				logger.Debugf("gcm: notification delivered: %s (%d)", gcmNotification.Identifier, userID)
			} else if len(result.Error) > 0 && result.Error == gcmNotRegisteredError {
				logger.Debugf("gcm: not registered: %s (%d)", gcmAuth.RegistrationID, userID)
			} else {
				logger.Errorf("gcm: notification couldn't be delivered: %s: %s (%d)", gcmNotification.Identifier, result.Error, userID)
			}
		}
	} else {
		logger.Errorf("gcm: status code %d", statusCode)
	}
	*/
}
