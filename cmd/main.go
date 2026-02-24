package main

import (
	"AIClinicReceptionist/internal/config"
	"AIClinicReceptionist/internal/database"
	"AIClinicReceptionist/internal/route"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()
	database.Connect(config.DBUrl)

	r := gin.Default()
	route.WebhookRouter(r)

	r.Run(":" + config.Port)
}
