package adm

import (
	"errors"
	"fmt"
	"os"

	"github.com/peterh/liner"
)

const (
	promptPassword = `Enter password: `
	repeatPassword = `Repeat password: `
	promptToken    = `Enter token: `
	repeatToken    = `Repeat token: `
	promptUser     = `Username: `

	SUCCESS = "\xe2\x9c\x94"
	FAILURE = "\xe2\x9c\x98"

	RED   = "\x1b[31m"
	GREEN = "\x1b[32m"
	CLEAR = "\x1b[0m"
)

var ErrPMismatch = errors.New(`Passwords did not match`)
var ErrTMismatch = errors.New(`Tokens did not match`)

func Read(style string) (string, error) {
	var (
		pass string
		err  error
	)

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	switch style {
	case `password`:
		pass, err = line.PasswordPrompt(promptPassword)
	case `token`:
		pass, err = line.PasswordPrompt(promptToken)
	case `user`:
		pass, err = line.Prompt(promptUser)
	}
	if err != nil {
		return "", err
	}
	return pass, nil
}

func ReadConfirmed(style string) (string, error) {
	var (
		pass, again string
		err         error
	)

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	// first pass
	switch style {
	case `password`:
		pass, err = line.PasswordPrompt(promptPassword)
	case `token`:
		pass, err = line.PasswordPrompt(promptToken)
	}
	if err != nil {
		return "", err
	}

	// second pass
	switch style {
	case `password`:
		again, err = line.PasswordPrompt(repeatPassword)
	case `token`:
		again, err = line.PasswordPrompt(repeatToken)
	}
	if err != nil {
		return "", err
	}

	// compare passwords
	if pass != again {
		switch style {
		case `password`:
			return "", ErrPMismatch
		case `token`:
			return "", ErrTMismatch
		}
	}

	return pass, nil
}

func ReadVerified(style string) string {
	var (
		password string
		err      error
	)

read_loop:
	for password == "" {
		if password, err = ReadConfirmed(style); err == liner.ErrPromptAborted {
			os.Exit(0)
		} else if err == ErrPMismatch || err == ErrTMismatch {
			fmt.Fprintf(
				os.Stderr,
				RED+FAILURE+CLEAR+" %s\n",
				err.Error(),
			)
			continue read_loop
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Fprintf(
		os.Stderr,
		GREEN+SUCCESS+CLEAR+" %s%s%s\n",
		`Entered `, style, `s match`,
	)
	return password
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
