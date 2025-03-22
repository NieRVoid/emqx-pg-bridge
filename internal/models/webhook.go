package models

import (
	"encoding/json"
)

// WebhookData represents the structure of the EMQX webhook data
type WebhookData struct {
	PublishReceivedAt int64                    `json:"publish_received_at"`
	Topic             string                   `json:"topic"`
	Payload           string                   `json:"payload"`
	Qos               int                      `json:"qos"`
	ClientID          string                   `json:"clientid"`
	Username          string                   `json:"username"`
	Timestamp         int64                    `json:"timestamp"`
	PubProps          WebhookPubProps          `json:"pub_props"`
	Event             string                   `json:"event"`
	Metadata          map[string]interface{}   `json:"metadata"`
	Flags             map[string]bool          `json:"flags"`
	Node              string                   `json:"node"`
	PeerHost          string                   `json:"peerhost"`
	PeerName          string                   `json:"peername"`
	ID                string                   `json:"id"`
	ClientAttrs       map[string]interface{}   `json:"client_attrs"`
}

// WebhookPubProps contains the MQTT publish properties
type WebhookPubProps struct {
	UserPropertyPairs []UserPropertyPair    `json:"User-Property-Pairs"`
	UserProperty      map[string]string     `json:"User-Property"`
}

// UserPropertyPair represents a key-value pair in the MQTT user properties
type UserPropertyPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetUserProperty returns a user property by key
func (d *WebhookData) GetUserProperty(key string) string {
	if d.PubProps.UserProperty != nil {
		if val, ok := d.PubProps.UserProperty[key]; ok {
			return val
		}
	}
	
	for _, prop := range d.PubProps.UserPropertyPairs {
		if prop.Key == key {
			return prop.Value
		}
	}
	
	return ""
}

// GetPayloadJSON parses the payload as JSON into the provided struct
func (d *WebhookData) GetPayloadJSON(v interface{}) error {
	return json.Unmarshal([]byte(d.Payload), v)
}