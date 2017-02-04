//
//  gcm.go
//  mercuryx
//
//  Copyright (c) 2016 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"fmt"
	"net/http"
	"github.com/mitchellh/mapstructure"
	"github.com/Hooks-Alerts/mercuryx/request"
	"github.com/Hooks-Alerts/mercuryx/config"
	"github.com/Hooks-Alerts/sirius/logger"
)

const GcmSenderID = "_android"
const ChromeSenderID = "chrome"

const gcmSendEndpoint = "https://android.googleapis.com/gcm/send"

const gcmNotRegisteredError = "NotRegistered"

func NewGcmSenderPool() *SenderPool {
	s := &SenderPool{ID: "gcm"}

	s.senderFactory = func() PushSender{
		return NewGcmPushSender()
	}
	s.initPool(config.Gcm.PoolSize)

	logger.Infof("gcm: initalized gcm sender (pool size: %d)", s.senderCount)
	return s
}

func NewChromeSenderPool() *SenderPool {
	s := &SenderPool{ID: "chrome"}

	s.senderFactory = func() PushSender{
		return NewGcmPushSender()
	}
	s.initPool(config.Chrome.PoolSize)

	logger.Infof("gcm: initalized chrome sender (pool size: %d)", s.senderCount)
	return s
}

// MARK: GcmPushSender

type GcmPushSender struct {
	client *http.Client
}

func NewGcmPushSender() *GcmPushSender {
	s := &GcmPushSender{}
	s.client = &http.Client{}
	return s
}

func (s *GcmPushSender) SendNotification(userID int, notification map[string]interface{}, auth map[string]interface{}) {
	gcmAuth := &GcmAuth{}
	if err := mapstructure.Decode(auth, gcmAuth); err != nil {
		logger.Errorf("gcm: invalid auth format: %v", err)
		return
	}

	gcmNotification := &GcmNotification{}
	if err := mapstructure.Decode(notification, gcmNotification); err != nil {
		logger.Errorf("gcm: invalid notification format: %v", err)
		return
	}

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
	r := request.NewRequest(s.client)
	defer r.Close()

	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("key=%s", gcmAuth.ApiKey)

	statusCode, err := r.URL(gcmSendEndpoint).Headers(headers).POST(&gcmReq).DoWithResponseEntity(&gcmResp)
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
}
