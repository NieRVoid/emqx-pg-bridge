package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	
	"github.com/NieRVoid/emqx-pg-bridge/internal/models"
)

// Common errors
var (
	ErrInvalidTopic = errors.New("invalid topic format")
)

// TopicInfo contains information extracted from a topic
type TopicInfo struct {
	RoomName   string
	DeviceName string
}

// ParseTopic parses a topic in the format "homestay/{roomName}/{deviceName}/"
// and extracts the room name and device name
func ParseTopic(topic string) (*TopicInfo, error) {
	parts := strings.Split(topic, "/")
	
	// The expected format is "homestay/{roomName}/{deviceName}/"
	// which splits into at least 3 parts
	if len(parts) < 3 || parts[0] != "homestay" {
		return nil, ErrInvalidTopic
	}
	
	return &TopicInfo{
		RoomName:   parts[1],
		DeviceName: parts[2],
	}, nil
}

// ParseRoomID extracts and validates the roomId from WebhookData
func ParseRoomID(data *models.WebhookData) (int, error) {
	roomIDStr := data.GetUserProperty("roomId")
	if roomIDStr == "" {
		return 0, errors.New("missing roomId in user properties")
	}
	
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid roomId format: %w", err)
	}
	
	return roomID, nil
}

// ParseDeviceID extracts and validates the deviceId from WebhookData
func ParseDeviceID(data *models.WebhookData) (int, error) {
	deviceIDStr := data.GetUserProperty("deviceId")
	if deviceIDStr == "" {
		return 0, errors.New("missing deviceId in user properties")
	}
	
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid deviceId format: %w", err)
	}
	
	return deviceID, nil
}

// ValidateJSON validates if a string is valid JSON
func ValidateJSON(data string) error {
	var js json.RawMessage
	return json.Unmarshal([]byte(data), &js)
}