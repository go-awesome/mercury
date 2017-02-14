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

type WebPushKeys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type WebPushSub struct {
	Endpoint string      `json:"endpoint"`
	Keys     WebPushKeys `json:"keys"`
}

type To struct {
	SenderID string      `json:"sender_id"`
	To       string      `json:"to,omitempty"`
	PushSub  *WebPushSub `json:"push_sub,omitempty"`
	Sandbox  bool        `json:"sandbox,omitempty"`
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
	Category         string      `json:"category,omitempty"`
	ContentAvailable bool        `json:"content_available,omitempty"`
	MutableContent   bool        `json:"mutable_content,omitempty"`
	Extra            interface{} `json:"extra,omitempty"`
}
