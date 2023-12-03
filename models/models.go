package models

import "time"

type MedicalRecord struct {
	Details               string
	DateOfRecord          time.Time // The original timestamp
	DateOfRecordFormatted string    // The formatted date as a string
}

type AccessControlRecord struct {
	PersonnelName      string    // The name of the medical personnel who has access
	StartDate          time.Time // The original timestamp
	StartDateFormatted string    // The formatted date as a string
}
