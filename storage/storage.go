//
//  storage.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package storage

import (
	"sync"
)

type storage interface {

	InsertSenderInfo(senderInfo *SenderInfo) error

	FetchSenderInfo(userID string, senderID string) (*SenderInfo, error)
	FetchSenderInfoArray(userID string) ([]*SenderInfo, error)

	DeleteSenderInfo(userID string, senderID string) error
}

// singleton interface
var  (
	instance storage
	once sync.Once
)

func Instance() storage {
	once.Do(func() {
		instance = NewMySql()
	})
	return instance
}
