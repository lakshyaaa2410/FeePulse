package controller

import (
	"fee-reminder/model"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (controller *Controller) AddMembersFromCSV(ctx *gin.Context) {

	// 1. Parse the uploaded CSV file
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please upload a CSV file with the form field name 'file'",
		})
		log.Printf("in controller.AddMembersFromCSV(): error parsing uploaded file: %v", err)
	}

	// 2. Validate file name and extension
	filename := file.Filename
	if !strings.HasSuffix(filename, ".csv") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid file type. Please upload a CSV file.",
		})
		log.Printf("in controller.AddMembersFromCSV(): invalid file type: %s", filename)
		return
	}

	// 3. Open the uploaded file
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error opening the uploaded file",
		})
		log.Printf("in controller.AddMembersFromCSV(): error opening uploaded file: %v", err)
		return
	}
	defer f.Close()

	// 4. Read the file content
	csvData, err := io.ReadAll(f)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading the uploaded file",
		})
		log.Printf("in controller.AddMembersFromCSV(): error reading uploaded file: %v", err)
		return
	}

	// 5. Process the CSV data and add members to the database
	err = controller.service.AddMembersFromCSV(csvData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error processing the CSV data",
		})
		log.Printf("in controller.AddMembersFromCSV(): error processing CSV data: %v", err)
		return
	}

	// 6. Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Members added successfully",
	})
}

func (controller *Controller) GetAllMembers(ctx *gin.Context) {
	members, err := controller.service.GetAllMembers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving members from the database",
		})
		log.Printf("in controller.GetAllMembers(): error retrieving members: %v", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": members,
	})
}

func (controller *Controller) AddNewMember(ctx *gin.Context) {

	var member model.Members

	err := ctx.BindJSON(&member)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body. Please provide member details in JSON format.",
		})
		log.Printf("in controller.AddNewMember(): error parsing request body: %v", err)
		return
	}

	// Validating required fields
	if member.Name == "" || member.Phone == "" || member.JoiningDate == "" || member.Duration <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing required fields. Please provide Name, Phone, JoiningDate, and Duration.",
		})
		log.Printf("in controller.AddNewMember(): missing required fields in request body: %+v", member)
		return
	}

	// Validating Phone Number
	if len(member.Phone) != 10 || !isNumeric(member.Phone) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid phone number. Please provide a 10-digit phone number.",
		})
		log.Printf("in controller.AddNewMember(): invalid phone number: %s", member.Phone)
		return
	}

	// Add the new member to the database
	err = controller.service.AddNewMember(member)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error adding new member to the database",
		})
		log.Printf("in controller.AddNewMember(): error adding member to database: %v", err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Member added successfully",
	})
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (controller *Controller) GetAllExpiringMemberships(ctx *gin.Context) {
	expiringMemberships, err := controller.service.GetAllExpiringMemberships()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error retrieving expiring memberships from the database",
		})
		log.Printf("in controller.GetAllExpiringMemberships(): error retrieving expiring memberships: %v", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": expiringMemberships,
	})
}
