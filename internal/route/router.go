package route

import (
	"AIClinicReceptionist/internal/database"
	"AIClinicReceptionist/internal/repository"
	"AIClinicReceptionist/internal/service"
	"AIClinicReceptionist/internal/webhooks"
	"time"

	"github.com/gin-gonic/gin"
)

func WebhookRouter(r *gin.Engine) {
	loc, _ := time.LoadLocation("Asia/Kolkata")

	doctorRepo := repository.NewDoctorRepository(database.DB)
	clinicRepo := repository.NewClinicHoursRepository(database.DB)
	appointmentRepo := repository.NewAppointmentRepository(database.DB)

	availabilityService := service.NewAvailabilityService(
		doctorRepo,
		clinicRepo,
		appointmentRepo,
		loc,
	)

	appointmentService := service.NewAppointmentService(doctorRepo, clinicRepo, appointmentRepo, loc)

	availabilityWebhook := webhooks.NewAvailabilityWebhook(availabilityService)
	appointmentWebhook := webhooks.NewAppointmentWebhook(appointmentService)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	Webhooks := r.Group("/webhooks")
	{
		Webhooks.POST("/book-appointment", appointmentWebhook.CreateAppointment)
		Webhooks.POST("/check-availability", availabilityWebhook.CheckAvailability)
	}
}
