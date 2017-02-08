//
//  sender.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import "github.com/ortuman/mercury/logger"

// MARK: PushSender

type SenderBuilder func() (PushSender, error)

type PushSender interface {
	SendNotification(userID int, notification map[string]interface{}, auth map[string]interface{})
}

// MARK: SenderPool

type SenderPool struct {
	ID				string
	senderCount 	int
	senderPool 		[]PushSender
	senderFactory 	func() PushSender
}

func (ph *SenderPool) SendNotification(userID int, notification map[string]interface{}, auth map[string]interface{}) {
	go func(userID int) {
		if ph.senderCount == 0 { return }
		ph.senderPool[userID % ph.senderCount].SendNotification(userID, notification, auth)
	}(userID)
}

func NewSenderPool(ID string, builder SenderBuilder, poolSize int) *SenderPool {
	sp := &SenderPool{}

	// assign pool identifier
	sp.ID = ID

	// initialize sender pool
	sp.senderCount = poolSize
	sp.senderPool  = make([]PushSender, 0, poolSize)

	for i := 0; i < sp.senderCount; i++ {
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
