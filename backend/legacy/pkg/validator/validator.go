package validator

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex = regexp.MustCompile(`^.{8,}$`)
)

// Init initializes custom validators
func Init() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validators
		v.RegisterValidation("password", validatePassword)
		v.RegisterValidation("safe_string", validateSafeString)
	}
	return nil
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	// At least 2 of 3 conditions
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	
	return count >= 2
}

// validateSafeString validates that string doesn't contain dangerous characters
func validateSafeString(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	// Check for common injection patterns
	dangerousPatterns := []string{
		"<script", "</script>", "javascript:", "onerror=", "onload=",
	}
	
	lowerStr := strings.ToLower(str)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerStr, pattern) {
			return false
		}
	}
	
	return true
}

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword validates a password
func ValidatePassword(password string) bool {
	return passwordRegex.MatchString(password)
}

// SanitizeString removes potentially dangerous content from a string
func SanitizeString(s string) string {
	// Remove HTML tags
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	s = tagRegex.ReplaceAllString(s, "")
	
	// Trim whitespace
	s = strings.TrimSpace(s)
	
	return s
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

