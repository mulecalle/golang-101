package factories

import (
	"errors"
	"golang-di-demo/contracts"
	drivers "golang-di-demo/services/email-drivers"
	"os"
	"strconv"
)

type EmailDriver int64

const (
	EmailServiceMessageBird EmailDriver = iota
	EmailServiceJustCall    EmailDriver = iota
)

func EmailDriverName(s string) (EmailDriver, error) {
	switch s {
	case "messagebird":
		return EmailServiceMessageBird, nil
	case "justcall":
		return EmailServiceJustCall, nil
	}
	return 0, errors.New("unable to parse the email driver name")
}

func EmailService(driver EmailDriver) (contracts.EmailInterface, error) {
	switch driver {
	case EmailServiceMessageBird:
		host := os.Getenv("MESSAGEBIRD_HOST")
		port, _ := strconv.Atoi(os.Getenv("MESSAGEBIRD_PORT"))

		return drivers.NewMessageBird(host, port), nil
	case EmailServiceJustCall:
		apiKey := os.Getenv("JUSTCALL_API_KEY")
		apiSecret := os.Getenv("JUSTCALL_API_SECRET")

		return drivers.NewJustCall(apiKey, apiSecret), nil
	}

	return nil, errors.New("Unable to create a new service")
}
