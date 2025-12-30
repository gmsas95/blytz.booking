package validator

import (
	"net/mail"
	"regexp"
	"unicode"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[\d\s\-\(\)]+$`)
	passwordRegex = regexp.MustCompile(`^.{6,}$`)
)

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email) || isEmailValid(email)
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidatePhone validates a phone number
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone) && len(phone) >= 10
}

// ValidatePassword validates a password
func ValidatePassword(password string) bool {
	if !passwordRegex.MatchString(password) {
		return false
	}

	// Check for at least one letter and one number
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}

// ValidateName validates a name
func ValidateName(name string) bool {
	return len(name) >= 2 && len(name) <= 100
}

// ValidateUUID validates a UUID string
func ValidateUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

// ValidateSlug validates a URL slug
func ValidateSlug(slug string) bool {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return slugRegex.MatchString(slug) && len(slug) >= 3 && len(slug) <= 50
}

// ValidateServiceName validates a service name
func ValidateServiceName(name string) bool {
	return len(name) >= 3 && len(name) <= 100
}

// ValidatePrice validates a price value
func ValidatePrice(price float64) bool {
	return price > 0 && price <= 100000
}
