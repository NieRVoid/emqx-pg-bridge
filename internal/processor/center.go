package processor

import (
	"context"
	"encoding/json"
	"strconv"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/NieRVoid/emqx-pg-bridge/internal/models"
	"github.com/NieRVoid/emqx-pg-bridge/pkg/logger"
)

// CenterProcessor handles processing for "device-center" type devices
type CenterProcessor struct {
	db  *pgxpool.Pool
	log *logger.Logger
}

// CenterPayload represents the payload structure for center devices
type CenterPayload struct {
	State           string `json:"state"`
	Count           int    `json:"count"`
	CountReliable   bool   `json:"count_reliable"`
	Source          string `json:"source"`
	Reliability     int    `json:"reliability"`
	Timestamp       int64  `json:"timestamp"`
}

// NewCenterProcessor creates a new center device processor
func NewCenterProcessor(db *pgxpool.Pool, log *logger.Logger) *CenterProcessor {
	return &CenterProcessor{
		db:  db,
		log: log,
	}
}

// Type returns the device type this processor handles
func (p *CenterProcessor) Type() string {
	return "device-center"
}

// Process handles the center device data
func (p *CenterProcessor) Process(ctx context.Context, data *models.WebhookData) error {
	// Extract roomId from user properties
	roomIDStr := data.GetUserProperty("roomId")
	if roomIDStr == "" {
		return ErrMissingRoomID
	}
	
	// Convert roomId to integer
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		p.log.Error("Invalid roomId format", "roomId", roomIDStr, "error", err)
		return err
	}
	
	// Parse payload
	var payload CenterPayload
	if err := json.Unmarshal([]byte(data.Payload), &payload); err != nil {
		p.log.Error("Failed to parse payload", "payload", data.Payload, "error", err)
		return ErrInvalidPayload
	}
	
	// Determine occupied status
	occupied := payload.State == "occupied"
	
	p.log.Debug("Processing center device data", 
		"roomId", roomID, 
		"occupied", occupied, 
		"count", payload.Count)
	
	// Update room_status table using UPSERT
	_, err = p.db.Exec(ctx, `
		INSERT INTO room_status (
			room_id, occupied, occupant_count, count_reliable, 
			count_source, source_reliability, updated_at, 
			last_occupancy_change
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), 
			CASE WHEN NOT EXISTS (
				SELECT 1 FROM room_status WHERE room_id = $1
			) OR $2 != (
				SELECT occupied FROM room_status WHERE room_id = $1
			) THEN 
				NOW() 
			ELSE 
				(SELECT last_occupancy_change FROM room_status WHERE room_id = $1)
			END
		)
		ON CONFLICT (room_id) 
		DO UPDATE SET
			occupied = $2,
			occupant_count = $3,
			count_reliable = $4,
			count_source = $5,
			source_reliability = $6,
			updated_at = NOW(),
			last_occupancy_change = CASE 
				WHEN room_status.occupied != $2 THEN NOW() 
				ELSE room_status.last_occupancy_change 
			END
	`, roomID, occupied, payload.Count, payload.CountReliable, 
	   payload.Source, payload.Reliability)
	
	if err != nil {
		p.log.Error("Failed to update room_status", "error", err)
		return err
	}
	
	p.log.Info("Updated room status", 
		"roomId", roomID, 
		"occupied", occupied, 
		"count", payload.Count)
	
	return nil
}