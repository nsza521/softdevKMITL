package utils

import (
	"errors"
	"strings"
)

func ToTitleCase(firstName, lastName string) (string, error) {

	if len(firstName) == 0 || len(lastName) == 0 {
		return "", errors.New("fullname must contain both of first and last name")
	}
	fullName := []string{firstName, lastName}
	for i, part := range fullName {
		if len(part) > 0 {
			fullName[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(fullName, " "), nil
}