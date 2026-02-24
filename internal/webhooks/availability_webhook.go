package webhooks

import (
	"AIClinicReceptionist/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CheckAvailabilityRequest struct {
	DoctorID string `json:"doctor_id" binding:"required"`
	Date     string `json:"date" binding:"required"` // YYYY-MM-DD
}

type CheckAvailabilityResponse struct {
	AvailableSlots []string `json:"available_slots"`
}

type AvailabilityWebhook struct {
	Service *service.AvailabilityService
}

func NewAvailabilityWebhook(s *service	.AvailabilityService) *AvailabilityWebhook {
	return &AvailabilityWebhook{Service: s}
}

func (h *AvailabilityWebhook) CheckAvailability(c *gin.Context) {

	var req CheckAvailabilityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	// Validate UUID
	doctorUUID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid doctor_id",
		})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date format (expected YYYY-MM-DD)",
		})
		return
	}

	slots, err := h.Service.GetAvailableSlots(
		c.Request.Context(),
		doctorUUID,
		date,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check availability",
		})
		return
	}

	c.JSON(http.StatusOK, CheckAvailabilityResponse{
		AvailableSlots: slots,
	})
}
