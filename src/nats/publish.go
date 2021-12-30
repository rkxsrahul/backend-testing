package nats

import (
	"log"

	"git.xenonstack.com/util/continuous-security-backend/config"
)

type RequestData struct {
	Method  string `json:"method"`
	UUID    string `json:"uuid"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Branch  string `json:"branch"`
	Subject string `json:"subject_type"`
}

func Publish(data []byte, subject string) {
	if err := config.NC.Publish(subject, data); err != nil {
		log.Fatal(err)
	}
}
