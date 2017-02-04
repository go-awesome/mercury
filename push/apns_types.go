package push

type ApnsAuth struct {
	Token	string	`mapstructure:"token"`
	Sandbox bool	`mapstructure:"sandbox"`
}

type ApnsNotification struct {
	Identifier	string	`mapstructure:"id,omitempty"`
	Title		string	`mapstructure:"alert_text,omitempty"`
	Body 		string	`mapstructure:"message,omitempty"`
	AlertID     string  `mapstructure:"alert_id,omitempty"`
	RoomID      string  `mapstructure:"room_id,omitempty"`
}

type ApnsRequest struct {
	APS APS	                    			`json:"aps,omitempty"`
	Notification   map[string]interface{} 	`json:"notification,omitempty"`
	NotificationID string       			`json:"id,omitempty"`
}

type APS struct {
	Alert APSAlert          `json:"alert,omitempty"`
	Badge uint              `json:"badge,omitempty"`
	Sound *string           `json:"sound,omitempty"`
	ContentAvailable int    `json:"content-available,omitempty"`
	MutableContent int      `json:"mutable-content,omitempty"`
	Category string         `json:"category,omitempty"`
}

type APSAlert struct {
	Title        string     `json:"title,omitempty"`
	TitleLocKey  string     `json:"title-loc-key,omitempty"`
	TitleLocArgs []string   `json:"title-loc-args,omitempty"`
	Body    string          `json:"body,omitempty"`
	LocKey  string          `json:"loc-key,omitempty"`
	LocArgs []string        `json:"loc-args,omitempty"`
	ActionLocKey string     `json:"action-loc-key,omitempty"`
	LaunchImage string      `json:"launch-image,omitempty"`
}
