//
//  push/apns_types.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

type ApnsRequest struct {
	APS          APS           `json:"aps,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
}

type APS struct {
	Alert            APSAlert `json:"alert,omitempty"`
	Badge            uint     `json:"badge,omitempty"`
	Sound            string   `json:"sound,omitempty"`
	ContentAvailable bool     `json:"content-available,omitempty"`
	MutableContent   bool     `json:"mutable-content,omitempty"`
	Category         string   `json:"category,omitempty"`
}

type APSAlert struct {
	Title        string   `json:"title,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
	Body         string   `json:"body,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
}
