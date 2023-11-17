package end2end

import "log"

type fakeLogger struct{}

// Log implements logging.Logger.
func (*fakeLogger) Logf(message string, a ...any) {
	log.Printf(message, a...)
}

func (*fakeLogger) Log(message string) {
	log.Println(message)
}
