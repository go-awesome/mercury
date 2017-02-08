//
//  sender.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"net/http"
	"crypto/tls"
	"golang.org/x/net/http2"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/cert"
	"github.com/ortuman/mercury/types"
)

const apnsSendEndpoint = "https://api.push.apple.com"
const apnsSandboxSendEndpoint = "https://api.development.push.apple.com"

func NewApnsSenderPool() *SenderPool {
	s := NewSenderPool("apns", NewApnsPushSender, config.Apns.PoolSize)
	logger.Infof("apns: initialized apns sender (pool size: %d)", s.senderCount)
	return s
}

// MARK: ApnsPushSender

type ApnsPushSender struct {
	client 			*http.Client
	sandboxClient 	*http.Client
}

func NewApnsPushSender() (PushSender, error) {
	s := &ApnsPushSender{}

	// Production certificate
	certificate, key, err := cert.LoadP12(config.Apns.CertFile, "")
	if err != nil {
		return nil, err
	}
	tlsCert := cert.TLSCertificate(certificate, key)

	// Sandbox connection
	sandboxCertificate, sandboxKey, err := cert.LoadP12(config.Apns.SandboxCertFile, "")
	if err != nil {
		return nil, err
	}
	tlsSandboxCert := cert.TLSCertificate(sandboxCertificate, sandboxKey)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{tlsCert}},
	}
	sandboxTransport := &http.Transport{
		TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{tlsSandboxCert}},
	}

	// Upgrade transports to HTTP2
	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
	}
	if err := http2.ConfigureTransport(sandboxTransport); err != nil {
		return nil, err
	}

	s.client = &http.Client{Transport: transport}
	s.sandboxClient = &http.Client{Transport: sandboxTransport}

	return s, nil
}

func (s *ApnsPushSender) SendNotification(userID string, notification *types.Notification) {

	// compose request
	/*
	apnsReq := ApnsRequest{}
	apnsReq.APS.Alert.Body = apnsNotification.Body
	apnsReq.APS.Alert.Title = apnsNotification.Title
	apnsReq.APS.Badge = 1 // badge

	soundFile := "base.mp3"
	apnsReq.APS.Sound = &soundFile

	apnsReq.APS.ContentAvailable = 1
	apnsReq.APS.MutableContent = 1
	apnsReq.APS.Category = "NOTIFICATION_CATEGORY"

	apnsReq.Notification = notification
	apnsReq.NotificationID = apnsNotification.Identifier

	var req *Request
	if !apnsAuth.Sandbox {
		req = NewRequest(s.client)
		req.URL(fmt.Sprintf("%v/3/device/%v", apnsSendEndpoint, apnsAuth.Token))
	} else {
		req = NewRequest(s.sandboxClient)
		req.URL(fmt.Sprintf("%v/3/device/%v", apnsSandboxSendEndpoint, apnsAuth.Token))
	}
	defer req.Close()

	// compose headers
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	expiration := time.Now().Add(86400 * time.Second)
	headers["apns-expiration"] = strconv.FormatInt(expiration.Unix(), 10)
	headers["apns-priority"] = "10"
	headers["apns-topic"] = "github.com/ortuman"

	statusCode, err := req.Headers(headers).POST().SendEntity(&apnsReq).Do(nil)
	if err != nil {
		logger.Errorf("apns: %v", err)
		return
	}

	var log string
	if statusCode == http.StatusOK {
		log = fmt.Sprintf("apns_sender: notification delivered: %s (%d)", apnsReq.NotificationID, apnsReq.APS.Badge)
	} else if statusCode == http.StatusGone {
		log = fmt.Sprintf("apns_sender: not registered: %s", apnsAuth.Token)
	} else {
		log = fmt.Sprintf("apns_sender: notification COULDN'T be delivered: %s (status: %v)", apnsReq.NotificationID, statusCode)
	}

	if apnsAuth.Sandbox {
		log += " [sandbox]"
	}
	logger.Debugf(log)
	*/
}
