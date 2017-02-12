//
//  gcm.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"fmt"
	"net/http"
	"golang.org/x/net/http2"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
	"encoding/json"
	"bytes"
)

const gcmSendEndpoint = "https://android.googleapis.com/gcm/send"

const gcmNotRegisteredError = "NotRegistered"

func NewGcmSenderPool() *SenderHub {
	s := NewSenderPool("gcm", NewGcmPushSender, config.Gcm.MaxConn)
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

func (s *GcmPushSender) SendNotification(to *To, notification *Notification) {

	gcmNotification := &GcmNotification{}
	gcmNotification.Title = notification.Title
	gcmNotification.Body = notification.Body
	gcmNotification.Sound = notification.Sound
	gcmNotification.Icon = notification.Icon
	gcmNotification.Color = notification.Color

	gcmReq := &GcmRequest{}
	gcmReq.RegistrationIDs = []string{to.To}
	gcmReq.Notification = gcmNotification
	gcmReq.Data = notification
	gcmReq.TimeToLive = 86400

	// perform request
	b, err := json.Marshal(&gcmReq)
	if err != nil {
		logger.Errorf("gcm: %v", err)
		return
	}

	req, err := http.NewRequest("POST", gcmSendEndpoint, bytes.NewReader(b))
	if err != nil {
		logger.Errorf("gcm: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", config.Gcm.ApiKey))

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Errorf("gcm: %v", err)
		return
	}
	defer resp.Body.Close()

	gcmResp := &GcmResponse{}
	if err := json.NewDecoder(resp.Body).Decode(gcmResp); err != nil {
		logger.Errorf("gcm: %v", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		for _, result := range gcmResp.Results {
			if len(result.Error) == 0 {
				logger.Debugf("gcm: notification delivered: %s (%s)", notification.ID, to.To)
			} else if len(result.Error) > 0 && result.Error == gcmNotRegisteredError {
				logger.Debugf("gcm: not registered: %s", to.To)
			} else {
				logger.Errorf("gcm: notification couldn't be delivered: %s: %s: (%s)", notification.ID, result.Error, to.To)
			}
		}
	} else {
		logger.Errorf("gcm: status code %d", resp.StatusCode)
	}
}
