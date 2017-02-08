//
//  types.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package storage

import "time"

type SenderInfo struct {
	UserID		string
	SenderID	string
	Token		string
	CreatedAt	time.Time
	UpdatedAt	time.Time
}
