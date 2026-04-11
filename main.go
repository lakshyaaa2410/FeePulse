// main.go
package main

import (
	"fee-reminder/controller"
	"fee-reminder/db"
	"fee-reminder/repository"
	"fee-reminder/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()

	db, err := db.InitializeDB()
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer db.Close()

	repository := repository.InitializeRepository(db)
	service := service.InitializeService(repository)
	controller := controller.InitializeController(service)

	// GET Methods
	router.GET("/members", controller.GetAllMembers)
	router.GET("/expiringMemberships", controller.GetAllExpiringMemberships)
	
	// POST Methods
	router.POST("/addMembersBulk", controller.AddMembersFromCSV)
	router.POST("/addMember", controller.AddNewMember)

	router.Run(":8080")
}

func init() {
	godotenv.Load()
}
