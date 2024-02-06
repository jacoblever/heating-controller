package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Logger interface {
	Logf(message string, a ...any)
	Log(message string)
}

type slackLogger struct {
	slackWebhookUrl string
}

func CreateSlackLogger() Logger {
	return &slackLogger{
		slackWebhookUrl: os.Getenv("SLACK_LOG_WEBHOOK_URL"),
	}
}

type SlackMessage struct {
	Text string `json:"text"`
}

func (l *slackLogger) Logf(message string, a ...any) {
	l.Log(fmt.Sprintf(message, a...))
}

// Log implements Logger.
func (l *slackLogger) Log(message string) {
	body := SlackMessage{
		Text: message,
	}
	jsonStr, err := json.Marshal(body)
	if err != nil {
		log.Printf("ERROR: failed to generate Slack request object: %s. Underlaying msg: %s", err, message)
	}

	req, err := http.NewRequest("POST", l.slackWebhookUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Printf("ERROR: failed to generate Slack log request: %s. Underlaying msg: %s", err, message)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: failed to send Slack log request: %s. Underlaying msg: %s", err, message)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: error from Slack sending log request, status code: %d. Underlaying msg: %s", resp.StatusCode, message)
	}
}
