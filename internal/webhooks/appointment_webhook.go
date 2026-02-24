package webhooks

import (
	"AIClinicReceptionist/internal/repository"
	"AIClinicReceptionist/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateAppointmentRequest struct {
	DoctorID    string `json:"doctor_id" binding:"required"`
	PatientName string `json:"patient_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	VisitType   string `json:"visit_type" binding:"required"`
	Date        string `json:"appointment_date" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
}

type AppointmentWebhook struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentWebhook(appointmentService *service.AppointmentService) *AppointmentWebhook {
	return &AppointmentWebhook{appointmentService: appointmentService}
}

func (h *AppointmentWebhook) CreateAppointment(c *gin.Context) {

	var req CreateAppointmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	doctorUUID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id"})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}

	err = h.appointmentService.CreateAppointment(
		c.Request.Context(),
		doctorUUID,
		req.PatientName,
		req.Phone,
		req.VisitType,
		date,
		req.StartTime,
	)

	if err != nil {
		if err == repository.ErrSlotAlreadyBooked {
			c.JSON(http.StatusConflict, gin.H{
				"error": "slot already booked",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "appointment booked successfully",
	})
}
