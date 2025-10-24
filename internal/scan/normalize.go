package scan

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// NormalizeE164 validates the supplied number and returns it formatted in E.164.
// If no defaultRegion is supplied, "US" is used to match PhoneInfoga defaults.
func NormalizeE164(number string, defaultRegion string) (string, error) {
	trimmed := strings.TrimSpace(number)
	if trimmed == "" {
		return "", fmt.Errorf("empty phone number")
	}
	if defaultRegion == "" {
		defaultRegion = "US"
	}
	parsed, err := phonenumbers.Parse(trimmed, defaultRegion)
	if err != nil {
		return "", fmt.Errorf("parse phone number: %w", err)
	}
	if !phonenumbers.IsValidNumber(parsed) {
		return "", fmt.Errorf("invalid phone number")
	}
	return phonenumbers.Format(parsed, phonenumbers.E164), nil
}
