package logging

import "github.com/jacoblever/heating-controller/brain/brain/clock"

type Logger interface {
	Logf(message string, a ...any)
	Log(message string)
}

type FileLogger interface {
	Logger
	DeleteAllLogs() error
}

type Loggers struct {
	clock       clock.Clock
	SlackLogger Logger
	Loggers     map[string]FileLogger
}

func InitLoggers(clock clock.Clock, slackLogger Logger) Loggers {
	return Loggers{
		clock:       clock,
		SlackLogger: slackLogger,
		Loggers:     map[string]FileLogger{},
	}
}

func (l *Loggers) Get(key string) Logger {
	logger, ok := l.Loggers[key]
	if !ok {
		l.SlackLogger.Logf("logger %s does not exist. Using SystemOutLogger instead", key)
		return &SystemOutLogger{}
	}
	return logger
}

func (l *Loggers) NewPerDayLogger(key string, settings Settings) {
	l.Loggers[key] = &PerDayLogger{
		Key:      key,
		Clock:    l.clock,
		Settings: settings,
	}
}

func (l *Loggers) CleanUpAnyLogs() {
	for _, v := range l.Loggers {
		v.DeleteAllLogs()
	}
}
