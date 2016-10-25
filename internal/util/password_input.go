package util

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/peterh/liner"
)

func (u *SomaUtil) GetNewPassword() string {
	var (
		newPassword string
		err         error
	)

	foundPassword := false
	for foundPassword == false {
		newPassword, err = u.ReadPasswordFromTerminal()
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
			continue
		}
		foundPassword = u.ValidatePasswordStrength(newPassword)

		if !foundPassword {
			fmt.Fprintf(os.Stdout, "Password entropy appears too low, try another one!\n\n")
		}
	}
	return newPassword
}

func (u *SomaUtil) ReadPasswordFromTerminal() (string, error) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	inpass01, err := line.PasswordPrompt("Enter new password: ")
	u.AbortOnError(err)

	inpass02, err := line.PasswordPrompt("Repeat new password: ")
	u.AbortOnError(err)

	if inpass01 == inpass02 {
		return inpass01, nil
	} else {
		err := errors.New("Input error: entered passwords did not match")
		return "", err
	}
}

func (u *SomaUtil) ValidatePasswordStrength(pass string) bool {
	matchLower := regexp.MustCompile(`[a-z]`)
	matchUpper := regexp.MustCompile(`[A-Z]`)
	matchNumber := regexp.MustCompile(`[0-9]`)
	matchSpecial := regexp.MustCompile(`[\!\@\#\$\%\^\&\*\(\)\-\\\_\=\+\,\.\?\/\:\;\{\}\[\]\~]`)
	matchFreaky := regexp.MustCompile(`[a-zA-Z0-9\!\@\#\$\%\^\&\*\(\)\-\\\_\=\+\,\.\?\/\:\;\{\}\[\]\~]`)

	score := 0

	// min length 8 chars
	if len(pass) < 8 {
		return false
	}

	// points based on length
	if len(pass) > 11 {
		score++
	}
	if len(pass) > 16 {
		score++
	}
	if len(pass) > 24 {
		score++
	}

	// points based on used character sets
	if matchLower.MatchString(pass) {
		score++
	}
	if matchUpper.MatchString(pass) {
		score++
	}
	if matchNumber.MatchString(pass) {
		score++
	}
	if matchSpecial.MatchString(pass) {
		score++
	}
	if matchFreaky.MatchString(pass) {
		score++
	}

	if score >= 4 {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
