package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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

// ValidateDate checks if the date string is valid based on specific rules.
func ValidateDate(date string) error {
	dateformat := "2006-01-02"
	//check if date is empty
	if date == "" {
		logrus.Debug("date is empty " + date)
		return errors.New("Date parameter is empty")
	}

	//check if date is in correct format
	parsedTime, err := time.Parse(dateformat, date)
	if err != nil {
		logrus.Debug("date is in correct format " + date + " parsedTime " + parsedTime.GoString())
		logrus.Error(err)
		return errors.New("Date is not in correct format")
	}

	//check if date is after current date
	currentTime := time.Now()
	if parsedTime.Before(currentTime) || parsedTime.Equal(currentTime) {
		logrus.Debug("date is after current date " + currentTime.GoString() + " parsedTime " + parsedTime.GoString())
		return errors.New("Date is before current date")
	}

	//check if the days correspond to the month

	// Extract the day from the parsed date
	day := parsedTime.Day()
	logrus.Debug("days correspond to the month - day ")
	logrus.Debug(day)
	// Extract the year and month from the parsed date
	year, month, _ := parsedTime.Date()
	logrus.Debug(month)
	logrus.Debug(year)
	// Calculate the last day of the month
	// Calculate the last day of the month
	lastDayOfMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour).Day()
	logrus.Debug(lastDayOfMonth)
	// Check if the day is valid for the chosen month
	if day < 1 && day > lastDayOfMonth {
		logrus.Debug("in the if 69")
		return errors.New("Date is out of range")
	}

	return nil
}
