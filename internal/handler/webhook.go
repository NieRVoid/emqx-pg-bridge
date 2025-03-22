package handler

import (
	"encoding/json"
	"io"
	"net/http"
	
	"github.com/NieRVoid/emqx-pg-bridge/internal/models"
	"github.com/NieRVoid/emqx-pg-bridge/internal/processor"
	"github.com/NieRVoid/emqx-pg-bridge/pkg/logger"
)

// WebhookHandler processes incoming webhooks from EMQX
type WebhookHandler struct {
	registry *processor.ProcessorRegistry
	log      *logger.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(registry *processor.ProcessorRegistry, log *logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		registry: registry,
		log:      log,
	}
}

// Handle processes webhook requests
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Limit request body size
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576)) // 1 MB limit
	if err != nil {
		h.log.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	
	// Parse the request body
	var data models.WebhookData
	if err := json.Unmarshal(body, &data); err != nil {
		h.log.Error("Failed to decode webhook data", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Extract device type
	deviceType := data.GetUserProperty("deviceType")
	if deviceType == "" {
		h.log.Error("Missing deviceType in webhook data")
		http.Error(w, "Missing deviceType", http.StatusBadRequest)
		return
	}
	
	h.log.Debug("Received webhook", 
		"deviceType", deviceType, 
		"topic", data.Topic, 
		"clientId", data.ClientID)
	
	// Get the appropriate processor
	processor, ok := h.registry.Get(deviceType)
	if !ok {
		h.log.Error("Unsupported device type", "deviceType", deviceType)
		http.Error(w, "Unsupported device type", http.StatusBadRequest)
		return
	}
	
	// Process the data
	if err := processor.Process(r.Context(), &data); err != nil {
		h.log.Error("Failed to process webhook data", 
			"deviceType", deviceType, 
			"error", err)
		http.Error(w, "Processing error", http.StatusInternalServerError)
		return
	}
	
	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}