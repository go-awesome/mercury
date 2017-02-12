//
//  types.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

const (
	ApnsSenderID = "apns"
	GcmSenderID = "gcm"
)

func IsValidSenderID(senderID string) bool {
	return senderID == ApnsSenderID || senderID == GcmSenderID
}

type To struct {
	SenderID string    `json:"sender_id"`
	To       string    `json:"to"`
	Sandbox  bool      `json:"sandbox,omitempty"`
}

type Push struct {
	To           []To         `json:"to"`
	Notification Notification `json:"notification"`
}

type Notification struct {
	ID               string      `json:"id"`
	Title            string      `json:"title,omitempty"`
	Body             string      `json:"body,omitempty"`
	Sound            string      `json:"sound,omitempty"`
	Icon             string      `json:"icon,omitempty"`
	Color            string      `json:"color,omitempty"`
	Category		 string		 `json:"category,omitempty"`
	ContentAvailable bool        `json:"content_available,omitempty"`
	MutableContent   bool        `json:"mutable_content,omitempty"`
	Extra            interface{} `json:"extra,omitempty"`
}
