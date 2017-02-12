//
//  sender.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"fmt"
	"time"
	"bytes"
	"strconv"
	"net/http"
	"crypto/tls"
	"encoding/json"
	"golang.org/x/net/http2"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/cert"
	"github.com/ortuman/mercury/storage"
)

const apnsSendEndpoint = "https://api.push.apple.com"
const apnsSandboxSendEndpoint = "https://api.development.push.apple.com"

func NewApnsSenderPool() *SenderHub {
	s := NewSenderPool("apns", NewApnsPushSender, config.Apns.MaxConn)
	logger.Infof("apns: initialized apns sender (max conn: %d)", s.senderCount)
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

func (s *ApnsPushSender) SendNotification(to *To, notification *Notification) {

	// compose request
	apnsReq := ApnsRequest{}
	apnsReq.APS.Alert.Body = notification.Body
	apnsReq.APS.Alert.Title = notification.Title

	badge, _ := storage.Instance().GetBadge(to.SenderID, to.To)
	apnsReq.APS.Badge = uint(badge)

	apnsReq.APS.Sound = notification.Sound
	apnsReq.APS.ContentAvailable = notification.ContentAvailable
	apnsReq.APS.MutableContent = notification.MutableContent
	apnsReq.APS.Category = notification.Category

	apnsReq.Notification = notification

	var sendEndpoint string
	if !to.Sandbox {
		sendEndpoint = apnsSendEndpoint
	} else {
		sendEndpoint = apnsSandboxSendEndpoint
	}

	// perform request
	b, err := json.Marshal(&apnsReq)
	if err != nil {
		logger.Errorf("apns: %v", err)
		return
	}

	req, err := http.NewRequest("POST", sendEndpoint, bytes.NewReader(b))
	if err != nil {
		logger.Errorf("apns: %v", err)
		return
	}
	expiration := time.Now().Add(86400 * time.Second)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apns-expiration", strconv.FormatInt(expiration.Unix(), 10))
	req.Header.Add("apns-priority", "10")
	req.Header.Add("apns-topic", "github.com/ortuman")

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Errorf("apns: %v", err)
		return
	}
	defer resp.Body.Close()

	var log string
	if resp.StatusCode == http.StatusOK {
		log = fmt.Sprintf("apns_sender: notification delivered: %s (%d)", notification.ID, apnsReq.APS.Badge)
	} else if resp.StatusCode == http.StatusGone {
		log = fmt.Sprintf("apns_sender: not registered: %s", to.To)
	} else {
		log = fmt.Sprintf("apns_sender: notification COULDN'T be delivered: %s (status: %v)", notification.ID, resp.StatusCode)
	}

	if to.Sandbox {
		log += " [sandbox]"
	}
	logger.Debugf(log)
}
