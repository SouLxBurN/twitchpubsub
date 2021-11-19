package schema

import "time"

type Message struct {
	Type string `json:"type,omitempty"`
	Data struct {
		Topic   string `json:"topic,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"data,omitempty"`
}

type TwitchPubMessage struct {
	Type string `json:"type"`
	Data struct {
		Timestamp  time.Time  `json:"timestamp,omitempty"`
		Redemption Redemption `json:"redemption,omitempty"`
	} `json:"data"`
}

type Redemption struct {
	ID        string     `json:"id"`
	User      TwitchUser `json:"user"`
	ChannelID string     `json:"channel_id"`
	Reward    Reward     `json:"reward"`
	UserInput string     `json:"user_input"`
	Status    string     `json:"status"`
}

type TwitchUser struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

type Reward struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Title     string `json:"title"`
}

type ListenRequestData struct {
	Topics    []string `json:"topics,omitempty"`
	AuthToken string   `json:"auth_token,omitempty"`
}

type ListenRequest struct {
	Type  string            `json:"type,omitempty"`
	Nonce string            `json:"nonce,omitempty"`
	Data  ListenRequestData `json:"data,omitempty"`
}
