package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"regexp"
	"strings"
	"unicode"
)

func generateSecureRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

var createRegexValidator = func(pattern string, emsg string) func(input string) error {
	return func(input string) error {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}

		if !regex.MatchString(input) {
			return errors.New(emsg)
		}

		return nil
	}
}

func isValidPassword(password string) bool {
	if len(password) < 10 {
		return false
	}

	var (
		hasLower   = false
		hasUpper   = false
		hasDigit   = false
		hasSpecial = false
	)

	specialChars := "!@#$%^&*()_+-=[]{};\\':\"|,./<>?"
	var lastChar rune
	var repeatCount int

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}

		if char == lastChar {
			repeatCount++
			if repeatCount == 3 {
				return false
			}
		} else {
			repeatCount = 1
		}
		lastChar = char
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}
