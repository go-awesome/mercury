//
//  push/gcm.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"golang.org/x/net/http2"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
	"time"
)

const gcmSendEndpoint = "https://android.googleapis.com/gcm/send"

const gcmNotRegisteredError = "NotRegistered"

func NewGcmSenderHub() *SenderHub {
	s := NewSenderHub(GcmSenderID, NewGcmPushSender, config.Gcm.MaxConn)
	logger.Infof("gcm: initialized gcm sender (max conn: %d)", s.senderCount)
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

func (s *GcmPushSender) SendNotification(to *To, notification *Notification) (int, time.Duration) {

	gcmNotification := &GcmNotification{}
	gcmNotification.Title = notification.Title
	gcmNotification.Body = notification.Body
	gcmNotification.Sound = notification.Sound
	gcmNotification.Icon = notification.Icon
	gcmNotification.Color = notification.Color

	gcmReq := &GcmRequest{}
	gcmReq.RegistrationIDs = []string{to.RegistrationID}
	gcmReq.Notification = gcmNotification
	gcmReq.Data = notification
	gcmReq.TimeToLive = 86400

	// perform request
	b, err := json.Marshal(&gcmReq)
	if err != nil {
		logger.Errorf("gcm: %v", err)
		return StatusFailed, 0
	}

	req, err := http.NewRequest("POST", gcmSendEndpoint, bytes.NewReader(b))
	if err != nil {
		logger.Errorf("gcm: %v", err)
		return StatusFailed, 0
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", config.Gcm.ApiKey))

	start := time.Now().UnixNano()
	resp, err := s.client.Do(req)
	end := time.Now().UnixNano()
	reqElapsed := time.Duration((end - start)) / time.Millisecond

	if err != nil {
		logger.Errorf("gcm: %v", err)
		return StatusFailed, reqElapsed
	}
	defer resp.Body.Close()

	gcmResp := &GcmResponse{}
	if err := json.NewDecoder(resp.Body).Decode(gcmResp); err != nil {
		logger.Errorf("gcm: %v", err)
		return StatusFailed, 0
	}

	if resp.StatusCode == http.StatusOK {
		for _, result := range gcmResp.Results {
			if len(result.Error) == 0 {
				logger.Debugf("gcm: [%s] notification delivered: %s (%s)", to.UserID, notification.ID, to.RegistrationID)
				return StatusDelivered, reqElapsed
			} else if len(result.Error) > 0 && result.Error == gcmNotRegisteredError {
				logger.Debugf("gcm: [%s] not registered: %s", to.UserID, to.RegistrationID)
				return StatusNotRegistered, 0
			} else {
				logger.Errorf("gcm: [%s] notification couldn't be delivered: %s: %s: (%s)", to.UserID, notification.ID, result.Error, to.RegistrationID)
				return StatusFailed, 0
			}
		}

	} else {
		logger.Errorf("gcm: status code %d", resp.StatusCode)
	}
	return StatusNone, 0
}
