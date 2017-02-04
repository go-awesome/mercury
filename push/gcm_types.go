//
//  gcm_types.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

type GcmAuth struct {
	ApiKey			string	`mapstructure:"api_key"`
	RegistrationID	string	`mapstructure:"registration_id"`
}

type GcmNotification struct {
	Identifier	string	`mapstructure:"id,omitempty"`
	Title		string	`mapstructure:"alert_text,omitempty"`
	Body 		string	`mapstructure:"message,omitempty"`
}

type GcmRequest struct {
	RegistrationIDs       []string               `json:"registration_ids"`
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
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
