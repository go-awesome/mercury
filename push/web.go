//
//  push/web.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/ortuman/mercury/config"
    "github.com/ortuman/mercury/logger"
    "github.com/ortuman/mercury/webpush"
    "golang.org/x/net/http2"
    "net/http"
    "strings"
    "time"
)

func NewChromeSenderHub() *SenderHub {
    s := NewSenderHub(ChromeSenderID, NewWebPushSender, config.Chrome.MaxConn)
    logger.Infof("web (%s): initialized chrome sender (max conn: %d)", ChromeSenderID, s.senderCount)
    return s
}

func NewFirefoxSenderHub() *SenderHub {
    s := NewSenderHub(FirefoxSenderID, NewWebPushSender, config.Firefox.MaxConn)
    logger.Infof("web (%s): initialized firefox sender (max conn: %d)", FirefoxSenderID, s.senderCount)
    return s
}

// MARK: WebPushSender

type WebPushSender struct {
    client *http.Client
}

func NewWebPushSender() (PushSender, error) {
    s := &WebPushSender{}
    transport := &http.Transport{}
    if err := http2.ConfigureTransport(transport); err != nil {
        return nil, err
    }
    s.client = &http.Client{Transport: transport}
    return s, nil
}

func (ws *WebPushSender) SendNotification(to *To, notification *Notification) (int, time.Duration) {
    b64 := base64.URLEncoding.WithPadding(base64.NoPadding)

    key, err := b64.DecodeString(strings.TrimRight(to.PushSub.Keys.P256dh, "="))
    if err != nil {
        logger.Errorf("web (%s): %v", to.SenderID, err)
        return StatusFailed, 0
    }
    auth, err := b64.DecodeString(strings.TrimRight(to.PushSub.Keys.Auth, "="))
    if err != nil {
        logger.Errorf("web (%s): %v", to.SenderID, err)
        return StatusFailed, 0
    }

    sub := &webpush.Subscription{Endpoint: to.PushSub.Endpoint, Key: key, Auth: auth}

    notificationJSON, err := json.Marshal(notification)
    if err != nil {
        logger.Errorf("web (%s): %v", to.SenderID, err)
        return StatusFailed, 0
    }

    wp := &webpush.Push{}

    subject, publicKey, privateKey, err := vapidDataFromSenderID(to.SenderID)
    if err != nil {
        logger.Errorf("web (%s): %v", to.SenderID, err)
        return StatusFailed, 0
    }

    wp.SetVapid(subject, publicKey, privateKey)

    start := time.Now().UnixNano()
    resp, err := wp.Do(ws.client, sub, string(notificationJSON), 86400)
    end := time.Now().UnixNano()
    reqElapsed := time.Duration((end - start)) / time.Millisecond

    if err != nil {
        logger.Errorf("web: %v", err)
        return StatusFailed, 0
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusCreated {
        logger.Debugf("web (%s): [%s] notification delivered: %s", to.SenderID, to.UserID, notification.ID)
        return StatusDelivered, reqElapsed

    } else if resp.StatusCode == http.StatusGone {
        logger.Debugf("web (%s): [%s] not registered: %s", to.SenderID, to.UserID, notification.ID)
        return StatusNotRegistered, 0

    } else {
        logger.Errorf("web (%s): [%s] notification couldn't be delivered: %s (status code: %d)", to.SenderID, to.UserID, notification.ID, resp.StatusCode)
        return StatusFailed, 0
    }
}

func vapidDataFromSenderID(senderID string) (string, string, string, error) {
    switch senderID {
    case ChromeSenderID:
        return config.Chrome.Subject, config.Chrome.PrivateKey, config.Chrome.PublicKey, nil

    case FirefoxSenderID:
        return config.Firefox.Subject, config.Firefox.PrivateKey, config.Firefox.PublicKey, nil
    }
    return "", "", "", errors.New(fmt.Sprintf("unrecognized web sender id: %s", senderID))
}
