package processor

import (
	"context"
	"errors"
	
	"github.com/NieRVoid/emqx-pg-bridge/internal/models"
	"github.com/NieRVoid/emqx-pg-bridge/pkg/logger"
)

// Common errors
var (
	ErrMissingDeviceID = errors.New("missing deviceId in user properties")
	ErrMissingRoomID   = errors.New("missing roomId in user properties")
	ErrInvalidPayload  = errors.New("invalid payload format")
)

// Processor defines the interface for all device type processors
type Processor interface {
	Process(ctx context.Context, data *models.WebhookData) error
	Type() string
}

// ProcessorRegistry maintains a mapping of device types to their processors
type ProcessorRegistry struct {
	processors map[string]Processor
	log        *logger.Logger
}

// NewProcessorRegistry creates a new registry
func NewProcessorRegistry(log *logger.Logger) *ProcessorRegistry {
	return &ProcessorRegistry{
		processors: make(map[string]Processor),
		log:        log,
	}
}

// Register adds a processor to the registry
func (r *ProcessorRegistry) Register(p Processor) {
	r.processors[p.Type()] = p
	r.log.Info("Registered processor", "type", p.Type())
}

// Get returns the processor for the given device type
func (r *ProcessorRegistry) Get(deviceType string) (Processor, bool) {
	p, ok := r.processors[deviceType]
	return p, ok
}

// GetProcessors returns all registered processors
func (r *ProcessorRegistry) GetProcessors() map[string]Processor {
	return r.processors
}