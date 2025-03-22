package processor

import (
	"context"
	"encoding/json"
	"strconv"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/NieRVoid/emqx-pg-bridge/internal/models"
	"github.com/NieRVoid/emqx-pg-bridge/pkg/logger"
)

// NormalProcessor handles processing for normal device types
type NormalProcessor struct {
	db  *pgxpool.Pool
	log *logger.Logger
}

// NewNormalProcessor creates a new normal device processor
func NewNormalProcessor(db *pgxpool.Pool, log *logger.Logger) *NormalProcessor {
	return &NormalProcessor{
		db:  db,
		log: log,
	}
}

// Type returns the device type this processor handles
func (p *NormalProcessor) Type() string {
	return "normal"
}

// Process handles the normal device data
func (p *NormalProcessor) Process(ctx context.Context, data *models.WebhookData) error {
	// Extract deviceId from user properties
	deviceIDStr := data.GetUserProperty("deviceId")
	if deviceIDStr == "" {
		return ErrMissingDeviceID
	}
	
	// Convert deviceId to integer
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		p.log.Error("Invalid deviceId format", "deviceId", deviceIDStr, "error", err)
		return err
	}
	
	p.log.Debug("Processing normal device data", "deviceId", deviceID)
	
	// For normal devices, we store the entire payload as JSON
	// Create a valid JSON object from the payload string
	var payloadJSON json.RawMessage
	if err := json.Unmarshal([]byte(data.Payload), &payloadJSON); err != nil {
		p.log.Error("Invalid JSON payload", "payload", data.Payload, "error", err)
		return ErrInvalidPayload
	}
	
	// Update device_status table using UPSERT
	_, err = p.db.Exec(ctx, `
		INSERT INTO device_status (
			device_id, status, updated_at, last_reported_at
		)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (device_id) 
		DO UPDATE SET
			status = $2,
			updated_at = NOW(),
			last_reported_at = NOW()
	`, deviceID, payloadJSON)
	
	if err != nil {
		p.log.Error("Failed to update device_status", "error", err)
		return err
	}
	
	p.log.Info("Updated device status", "deviceId", deviceID)
	
	return nil
}