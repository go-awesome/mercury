//
//  push/gcm_types.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

type GcmNotification struct {
	Title        string    `json:"title,omitempty"`
	Body         string    `json:"body,omitempty"`
	Icon         string    `json:"icon,omitempty"`
	Sound        string    `json:"sound,omitempty"`
	Tag          string    `json:"tag,omitempty"`
	Color        string    `json:"color,omitempty"`
	ClickAction  string    `json:"click_action,omitempty"`
	BodyLocKey   string    `json:"body_loc_key,omitempty"`
	BodyLocArgs  string    `json:"body_loc_args,omitempty"`
	TitleLocKey  string    `json:"title_loc_key,omitempty"`
	TitleLocArgs string    `json:"title_loc_args,omitempty"`
}

type GcmRequest struct {
	RegistrationIDs       []string         `json:"registration_ids"`
	CollapseKey           string           `json:"collapse_key,omitempty"`
	DelayWhileIdle        bool             `json:"delay_while_idle,omitempty"`
	TimeToLive            int              `json:"time_to_live,omitempty"`
	RestrictedPackageName string           `json:"restricted_package_name,omitempty"`
	DryRun                bool             `json:"dry_run,omitempty"`
	Notification          *GcmNotification `json:"notification,omitempty"`
	Data                  interface{}       `json:"data,omitempty"`
}

type GcmResponse struct {
	MulticastID  int64       `json:"multicast_id"`
	Success      int         `json:"success"`
	Failure      int         `json:"failure"`
	CanonicalIDs int         `json:"canonical_ids"`
	Results      []GcmResult `json:"results"`
}

// Result represents the status of a processed message.
type GcmResult struct {
	MessageID      string `json:"message_id"`
	RegistrationID string `json:"registration_id"`
	Error          string `json:"error"`
}
