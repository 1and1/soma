package adm

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/nbutton23/zxcvbn-go"
	"github.com/nbutton23/zxcvbn-go/scoring"
	"github.com/peterh/liner"
)

func EvaluatePassword(minScore int, sl ...string) (bool, error) {
	var (
		err     error
		quality scoring.MinEntropyMatch
		choice  string
	)

	if len(sl) < 1 {
		return false, fmt.Errorf(`No password given for judgement`)
	}

	// workaround https://github.com/nbutton23/zxcvbn-go/issues/15
	if _, err = hex.DecodeString(sl[0]); err == nil {
		if len(sl[0]) >= 16 {
			return true, nil
		}
		return false, fmt.Errorf(`Hexadecimal strings must be 16 characters or more`)
	}

	// second argument are additional punished strings
	quality = zxcvbn.PasswordStrength(sl[0], sl[1:])

	// display evaluation summary report
	fmt.Printf(
		`Password score    (0-4): %d
Estimated entropy (bit): %f
Estimated time to crack: %s%s`,
		quality.Score,
		quality.Entropy,
		quality.CrackTimeDisplay, "\n",
	)

	// enfore chance for a better password
	if quality.Score < minScore {
		fmt.Println(RED + FAILURE + CLEAR +
			` Chosen password is too weak. Please select a better one.`)
		return false, nil
	}

	// offer chance for a better password
	for choice != `y` && choice != `n` {
		if choice, err = verifyChoice(); err != nil {
			if err == liner.ErrPromptAborted {
				os.Exit(0)
			}
			return false, err
		}
	}

	switch choice {
	case `y`:
		return true, nil
	case `n`:
		return false, nil
	}
	return false, fmt.Errorf(`Unreachable error reached`)
}

func verifyChoice() (string, error) {
	var (
		err    error
		choice string
	)

	const (
		prompt string = `Select this password? (y/n): `
	)
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	if choice, err = line.Prompt(prompt); err != nil {
		return "", err
	}

	return choice, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
