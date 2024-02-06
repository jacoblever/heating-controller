package logging

import "log"

type SystemOutLogger struct{}

func (*SystemOutLogger) Logf(message string, a ...any) {
	log.Printf(message, a...)
}

func (*SystemOutLogger) Log(message string) {
	log.Println(message)
}
