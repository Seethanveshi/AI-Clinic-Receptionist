package main

import (
	"AIClinicReceptionist/internal/config"
	"AIClinicReceptionist/internal/database"
)

func main() {
	config := config.Load()
	database.Connect(config.DBUrl)

	// r := gin.Default()

}
