package models

import "time"

type MedicalRecord struct {
	Details               string
	DateOfRecord          time.Time // The original timestamp
	DateOfRecordFormatted string    // The formatted date as a string
}
