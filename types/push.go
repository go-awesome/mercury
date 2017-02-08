//
//  push.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package types

type Push struct {
	SenderID 	 string         `json:"sender_id,omitempty"`
	UserID 		 string			`json:"user_id"`
	Notification *Notification	`json:"notification"`
}

type Notification struct {
	ID	string	`json:"id"`
}
