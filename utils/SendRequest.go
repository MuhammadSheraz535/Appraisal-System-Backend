package utils

import (
	"bytes"
	"net/http"
	"os"
	"time"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
)

func SendRequest(method string, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	token := os.Getenv("TOSS_BEARER_TOKEN")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: time.Second * 20,
	}

	log.Info("Sending " + method + " request to " + url)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return resp, nil
}
