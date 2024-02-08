package logging

import "log"

type SystemOutLogger struct{}

func (*SystemOutLogger) Logf(message string, a ...any) {
	log.Printf(message, a...)
}

func (*SystemOutLogger) Log(message string) {
	log.Println(message)
}

func (*SystemOutLogger) GetLogs() ([]string, error) {
	return []string{}, nil
}

func (*SystemOutLogger) DeleteAllLogs() error {
	return nil
}
