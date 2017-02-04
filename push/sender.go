//
//  sender.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

// MARK: PushSender

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

func (ph *SenderPool) initPool(poolSize int) {

	// initialize sender pool
	ph.senderCount = poolSize
	ph.senderPool  = make([]PushSender, 0, poolSize)

	for i := 0; i < ph.senderCount; i++ {
		ps := ph.senderFactory()
		if ps == nil { ph.senderCount = 0; break }
		ph.senderPool = append(ph.senderPool, ps)
	}
}
