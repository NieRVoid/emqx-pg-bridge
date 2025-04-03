package processor

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/NieRVoid/emqx-pg-bridge/internal/models"
	"github.com/NieRVoid/emqx-pg-bridge/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CenterProcessor handles processing for "device-center" type devices
type CenterProcessor struct {
	db  *pgxpool.Pool
	log *logger.Logger
}

// CenterPayload represents the payload structure for center devices
type CenterPayload struct {
	Count              int    `json:"count"`
	CountConfidence    int    `json:"countConfidence"`
	OccupiedConfidence int    `json:"occupiedConfidence"`
	LastResetTime      int    `json:"lastResetTime"`
	LastChangeTime     int64  `json:"lastChangeTime"`
	ChangeSource       string `json:"changeSource"`
	Occupied           bool   `json:"occupied"`
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

	p.log.Debug("Processing center device data",
		"roomId", roomID,
		"occupied", payload.Occupied,
		"count", payload.Count,
		"countConfidence", payload.CountConfidence,
		"occupiedConfidence", payload.OccupiedConfidence)

	// Update room_status table using UPSERT
	_, err = p.db.Exec(ctx, `
		INSERT INTO room_status (
			room_id, occupied, occupant_count, count_confidence, 
			occupied_confidence, count_source, updated_at, 
			last_source_change
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), 
			CASE WHEN NOT EXISTS (
				SELECT 1 FROM room_status WHERE room_id = $1
			) OR $6 != (
				SELECT count_source FROM room_status WHERE room_id = $1
			) THEN 
				NOW() 
			ELSE 
				(SELECT last_source_change FROM room_status WHERE room_id = $1)
			END
		)
		ON CONFLICT (room_id) 
		DO UPDATE SET
			occupied = $2,
			occupant_count = $3,
			count_confidence = $4,
			occupied_confidence = $5,
			count_source = $6,
			updated_at = NOW(),
			last_source_change = CASE 
				WHEN room_status.count_source != $6 THEN NOW() 
				ELSE room_status.last_source_change 
			END
	`, roomID, payload.Occupied, payload.Count, payload.CountConfidence,
		payload.OccupiedConfidence, payload.ChangeSource)

	if err != nil {
		p.log.Error("Failed to update room_status", "error", err)
		return err
	}

	p.log.Info("Updated room status",
		"roomId", roomID,
		"occupied", payload.Occupied,
		"count", payload.Count,
		"countConfidence", payload.CountConfidence,
		"occupiedConfidence", payload.OccupiedConfidence)

	return nil
}
