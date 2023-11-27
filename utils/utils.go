package utils

import (
	"errors"
	"strconv"
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
