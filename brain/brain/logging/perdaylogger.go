package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/fileio"
	"github.com/jacoblever/heating-controller/brain/brain/timeseries"
)

type Settings struct {
	DaysToKeepFor       int
	SynchronousDeletion bool
}

type PerDayLogger struct {
	Key      string
	Clock    clock.Clock
	Settings Settings
}

func (l *PerDayLogger) Logf(message string, a ...any) {
	l.Log(fmt.Sprintf(message, a...))
}

func (l *PerDayLogger) Log(message string) {
	err := l.append(message)
	if err != nil {
		log.Println("[PerDayLogger.append] error appending: ", err)
	}
}

func (l *PerDayLogger) GetLogs() ([]string, error) {
	logs := []string{}
	entries, err := os.ReadDir(l.Key)
	if err != nil {
		return logs, err
	}

	for _, e := range entries {
		_, err := time.Parse("2006-01-02.txt", e.Name())
		if err != nil {
			continue
		}

		data, err := fileio.ReadFile(path.Join(l.Key, e.Name()))
		if err != nil {
			return logs, err
		}

		lines := strings.Split(data, "\n")
		lines = lines[1:]
		logs = append(logs, lines...)
	}

	return logs, nil
}

func (l *PerDayLogger) DeleteAllLogs() error {
	return os.RemoveAll(l.Key)
}

func (l PerDayLogger) append(message string) error {
	now := l.Clock.Now()
	dayFileName := fmt.Sprintf("%s/%02d-%02d-%02d.txt", l.Key, now.Year(), now.Month(), now.Day())

	if err := os.MkdirAll(l.Key, 0770); err != nil {
		return err
	}

	err := timeseries.Append(dayFileName, l.Clock, message, nil)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := l.deleteOldFiles()
		if err != nil {
			log.Println("[PerDayLogger.append] error deleting old files: ", err)
		}
		wg.Done()
	}()

	if l.Settings.SynchronousDeletion {
		wg.Wait()
	}
	return nil
}

func (l PerDayLogger) deleteOldFiles() error {
	entries, err := os.ReadDir(l.Key)
	if err != nil {
		return err
	}

	for _, e := range entries {
		date, err := time.Parse("2006-01-02.txt", e.Name())
		if err != nil {
			continue
		}

		if date.Add(time.Duration(l.Settings.DaysToKeepFor) * 24 * time.Hour).Before(l.Clock.Now()) {
			err := os.Remove(path.Join(l.Key, e.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
