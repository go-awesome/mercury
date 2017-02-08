//
//  push.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package types

type Push struct {
	SenderID 	 string         `json:"sender_id,omitempty"`
	UserIDs 	 []string		`json:"user_ids"`
	Notification Notification	`json:"notification"`
}

type Notification struct {
	ID	string	`json:"id"`
}

