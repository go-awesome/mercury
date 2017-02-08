//
//  sender.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"hash/fnv"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/types"
)

// MARK: PushSender

type PushSender interface {
	SendNotification(userID string, notification *types.Notification)
}

// MARK: SenderPool

type SenderPool struct {
	ID			string
	senderCount	uint32
	senderPool 	[]PushSender
}

func NewSenderPool(ID string, builder func() (PushSender, error), poolSize uint32) *SenderPool {
	sp := &SenderPool{}

	// assign pool identifier
	sp.ID = ID

	// initialize sender pool
	sp.senderCount = poolSize
	sp.senderPool  = make([]PushSender, 0, poolSize)

	for i := uint32(0); i < sp.senderCount; i++ {
		ps, err := builder()
		if err != nil {
			logger.Errorf("sender: %v", err)
			sp.senderCount = 0
			break
		}
		sp.senderPool = append(sp.senderPool, ps)
	}
	return sp
}

func (ph *SenderPool) SendNotification(userID string, notification *types.Notification) {
	go ph.send(userID, notification)
}

func (ph *SenderPool) send(userID string, notification *types.Notification) {
	h := fnv.New32a()
	h.Write([]byte(userID))
	ph.senderPool[h.Sum32() % ph.senderCount].SendNotification(userID, notification)
}
