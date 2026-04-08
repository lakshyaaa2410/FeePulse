package service

import (
	"bytes"
	"encoding/csv"
	"fee-reminder/model"
	"fmt"
	"strconv"
	"time"
)

func (service *Service) AddMembersFromCSV(csvData []byte) error {

	// 1. Creating a new reader for the CSV data
	reader := csv.NewReader(bytes.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("in service.AddMembersFromCSV(): error parsing CSV data: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("in service.AddMembersFromCSV(): CSV file must contain at least a header and one data row")
	}

	if len(records[0]) < 4 {
		return fmt.Errorf("in service.AddMembersFromCSV(): CSV file must contain at least 4 columns: Name, Joining Date, Phone, Duration")
	}

	var members []model.Members

	// 2. Process each record and add members to the database
	for _, row := range records[1:] {
		name := row[0]
		phone := row[1]
		joiningDate := row[2]

		// Validate joining date format (dd-mm-yyyy)
		duration, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			fmt.Printf("in service.AddMembersFromCSV(): error parsing duration for member %s: %v\n", name, err)
			continue
		}

		// Calculate expiry date based on joining date and duration
		expiryDate := calculateExpiryDate(joiningDate, duration)

		members = append(members, model.Members{
			Name:        name,
			Phone:       phone,
			JoiningDate: joiningDate,
			Duration:    duration,
			ExpiryDate:  expiryDate,
		})
	}

	// 3. Add members to the database
	err = service.repository.AddMembers(members)
	if err != nil {
		return fmt.Errorf("in service.AddMembersFromCSV(): error adding members to database: %v", err)
	}

	// 4. Return success response
	return nil
}

func calculateExpiryDate(joiningDate string, duration int64) int64 {
	ist, _ := time.LoadLocation("Asia/Kolkata")

	joiningDateTime, err := time.ParseInLocation("02-01-2006", joiningDate, ist)
	if err != nil {
		fmt.Printf("in service.calculateExpiryDate(): error parsing joining date: %v\n", err)
		return 0
	}

	expiryDateTime := joiningDateTime.AddDate(0, int(duration), 0)
	return expiryDateTime.Unix()
}
