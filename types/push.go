//
//  push.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package types

type Push struct {
	SenderID 	 string         		`json:"sender_id"`
	UserID 		 int					`json:"user_id"`
	Auth 		 map[string]interface{} `json:"auth"`
	Notification map[string]interface{}	`json:"notification"`
}

