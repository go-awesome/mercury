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
)

// MARK: PushSender

type PushSender interface {
	SendNotification(to *To, notification *Notification)
}

// MARK: SenderHub

type SenderHub struct {
	ID			string
	senderCount	uint32
	senderPool 	[]PushSender
}

func NewSenderPool(ID string, builder func() (PushSender, error), poolSize uint32) *SenderHub {
	sh := &SenderHub{}

	// assign pool identifier
	sh.ID = ID

	// initialize sender pool
	sh.senderCount = poolSize
	sh.senderPool  = make([]PushSender, 0, poolSize)

	for i := uint32(0); i < sh.senderCount; i++ {
		ps, err := builder()
		if err != nil {
			logger.Errorf("sender: %v", err)
			sh.senderCount = 0
			break
		}
		sh.senderPool = append(sh.senderPool, ps)
	}
	return sh
}

func (sh *SenderHub) SendNotification(to *To, notification *Notification) {
	go sh.send(to, notification)
}

func (sh *SenderHub) send(to *To, notification *Notification) {
	h := fnv.New32a()
	h.Write([]byte(to.SenderID + ":" + to.To))
	sh.senderPool[h.Sum32() % sh.senderCount].SendNotification(to, notification)
}
