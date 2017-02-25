//
//  storage/storage.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package storage

import (
	"sync"
)

type storage interface {
	IncreaseBadge(senderID, token string) error
	GetBadge(senderID, token string) (uint64, error)
	ClearBadge(senderID, token string) error
}

// singleton interface
var  (
	instance storage
	once sync.Once
)

func Instance() storage {
	once.Do(func() {
		instance = NewRedis()
	})
	return instance
}
