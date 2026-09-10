package idp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// ParseInt converts a string to int. Panics on invalid input.
func ParseInt(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

// GetCountryCode returns the ISO country code for the given phone prefix and number.
func GetCountryCode(prefix string, phone string) (string, error) {
	if prefix == "" || phone == "" {
		return "", nil
	}

	phoneNumber, err := phonenumbers.Parse(fmt.Sprintf("+%s%s", prefix, phone), "")
	if err != nil {
		return "", err
	}

	countryCode := phonenumbers.GetRegionCodeForNumber(phoneNumber)
	if countryCode == "" {
		return "", fmt.Errorf("country code not found for phone prefix: %s", prefix)
	}

	return countryCode, nil
}

// GetUsernameFromEmail extracts the username part before '@' from an email address.
func GetUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	return parts[0]
}

// IsChinese checks whether the given string contains any Chinese characters.
func IsChinese(str string) bool {
	for _, r := range str {
		if r >= '\u4e00' && r <= '\u9fa5' {
			return true
		}
	}
	return false
}

// ParseIdToString converts various ID types to string.
func ParseIdToString(input interface{}) (string, error) {
	switch v := input.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported id type: %T", input)
	}
}
