//
//  push/sender_hub.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"errors"
	"fmt"
	"github.com/ortuman/mercury/config"
	"github.com/ortuman/mercury/logger"
	"golang.org/x/net/http2"
	"hash/fnv"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// MARK: PushSender

const (
	StatusNone = iota
	StatusDelivered
	StatusNotRegistered
	StatusFailed
)

type PushSender interface {
	SendNotification(to *To, notification *Notification) (int, time.Duration)
}

type PushStats struct {
	MaxConn           uint32 `json:"max_conn"`
	DeliveredCount    uint64 `json:"delivered_count"`
	UnregisteredCount uint64 `json:"unregistered_count"`
	FailedCount       uint64 `json:"failed_count"`
	AvgRequestTime    uint64 `json:"avg_request_time_msec"`
}

// MARK: SenderHub

type SenderHub struct {
	ID                string
	senderCount       uint32
	senderPool        []PushSender
	deliveredCount    uint64
	unregisteredCount uint64
	failedCount       uint64
	sumRequestTime    uint64
}

var unregisteredCallbackClient *http.Client

func init() {
	transport := &http.Transport{}
	if err := http2.ConfigureTransport(transport); err != nil {
		log.Fatalf("sender_hub: %v", err)
	}
	unregisteredCallbackClient = &http.Client{Transport: transport}
}

func NewSenderHub(ID string, builder func() (PushSender, error), poolSize uint32) *SenderHub {
	sh := &SenderHub{}

	// assign hub identifier
	sh.ID = ID

	// initialize sender hub
	sh.senderCount = poolSize
	sh.senderPool = make([]PushSender, 0, poolSize)

	for i := uint32(0); i < sh.senderCount; i++ {
		ps, err := builder()
		if err != nil {
			logger.Errorf("sender_hub: %v", err)
			sh.senderCount = 0
			break
		}
		sh.senderPool = append(sh.senderPool, ps)
	}
	return sh
}

func (sh *SenderHub) GetID() string {
	return sh.ID
}

func (sh *SenderHub) SendNotification(to *To, notification *Notification) int {
	go sh.send(to, notification)
	return StatusNone
}

func (sh *SenderHub) Stats() PushStats {
	stats := PushStats{
		MaxConn:           sh.senderCount,
		DeliveredCount:    atomic.LoadUint64(&sh.deliveredCount),
		UnregisteredCount: atomic.LoadUint64(&sh.unregisteredCount),
		FailedCount:       atomic.LoadUint64(&sh.failedCount),
	}
	if stats.DeliveredCount > 0 {
		stats.AvgRequestTime = atomic.LoadUint64(&sh.sumRequestTime) / stats.DeliveredCount
	}
	return stats
}

func (sh *SenderHub) send(to *To, notification *Notification) {
	h := fnv.New32a()
	h.Write([]byte(to.UserID))

	status, reqElapsed := sh.senderPool[h.Sum32()%sh.senderCount].SendNotification(to, notification)

	switch status {
	case StatusDelivered:
		atomic.AddUint64(&sh.deliveredCount, 1)
		atomic.AddUint64(&sh.sumRequestTime, uint64(reqElapsed))

	case StatusNotRegistered:
		if err := notifyUnregistered(to.SenderID, to.UserID); err != nil {
			logger.Errorf("sender_hub: %v", err)
		} else {
			atomic.AddUint64(&sh.unregisteredCount, 1)
		}

	case StatusFailed:
		atomic.AddUint64(&sh.failedCount, 1)

	default:
		break
	}
}

func notifyUnregistered(senderID, userID string) error {
	unregisteredURL := config.Server.UnregisteredCallback + "/" + senderID + "/" + userID

	req, err := http.NewRequest("GET", unregisteredURL, nil)
	if err != nil {
		return err
	}

	resp, err := unregisteredCallbackClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("sender_hub: unregistered callback: status code %d", resp.StatusCode))
	}
	return nil
}
