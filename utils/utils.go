package utils

import (
	"errors"
	"strconv"
	"time"
)

// ValidateNumeroUtente checks if the numeroUtente string is valid based on specific rules.
func ValidateNumeroUtente(numeroUtente string) error {
	// Check if the string is empty
	if numeroUtente == "" {
		return errors.New("Numero utente parameter is missing")
	}

	// Check if the length is exactly 9 digits
	if len(numeroUtente) != 9 {
		return errors.New("Numero utente must be exactly 9 digits long")
	}

	// Check if all characters are numbers
	if _, err := strconv.Atoi(numeroUtente); err != nil {
		return errors.New("Numero utente must contain only numbers")
	}

	return nil
}

func ValidateDate(date string) error {
	dateformat := "02/01/2006"
	//check if date is empty
	if date == "" {
		return errors.New("Date parameter is empty")
	}

	//check if date is in correct format
	parsedTime, err := time.Parse(dateformat, date)
	if err != nil {
		return errors.New("Date is not in correct format")
	}

	//check if date is after current date
	currentTime := time.Now()
	if parsedTime.Before(currentTime) || parsedTime.Equal(currentTime) {
		return errors.New("Date is before current date")
	}

	//check if the days correspond to the month

	// Extract the day from the parsed date
	day := parsedTime.Day()

	// Extract the year and month from the parsed date
	year, month, _ := parsedTime.Date()

	// Calculate the last day of the month
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	// Check if the day is valid for the chosen month
	if day >= 1 && day <= lastDayOfMonth {
		return errors.New("Date is out of range")
	}

	return nil
}
