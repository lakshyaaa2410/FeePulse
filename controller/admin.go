package controller

import (
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
