package models

import "time"

type MedicalRecord struct {
	Description string `json:"description"`
	CreatedDate int64  `json:"createDate"`
	Date        int64  `json:"date"`
	EntityName  string `json:"entityName"`
	RecordType  string `json:"type"`
}

type AccessControlRecord struct {
	PersonnelName      string    // The name of the medical personnel who has access
	StartDate          time.Time // The original timestamp
	StartDateFormatted string    // The formatted date as a string
	AccessStatus       string    // The access status
}
