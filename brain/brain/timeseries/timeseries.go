package timeseries

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jacoblever/heating-controller/brain/brain/clock"
	"github.com/jacoblever/heating-controller/brain/brain/fileio"
)

func Append(filePath string, clock clock.Clock, value string, onlyEvery *time.Duration) error {
	if shouldAppend(filePath, clock, onlyEvery) {
		timestamp := clock.Now().Format(time.RFC3339)
		line := strings.Join([]string{timestamp, value}, ",")
		return fileio.AppendLineToFile(filePath, line)
	}
	return nil
}

func ReadLastRecord(filePath string) ([]string, error) {
	line, err := fileio.ReadLastLine(filePath)
	if err != nil {
		return nil, err
	}

	if line == "" {
		return nil, fmt.Errorf("no last record in %s", filePath)
	}
	parts := strings.Split(line, ",")
	return parts[1:], nil
}

func shouldAppend(filePath string, clock clock.Clock, onlyEvery *time.Duration) bool {
	if onlyEvery == nil {
		return true
	}

	lastLine, err := fileio.ReadLastLine(filePath)
	if err != nil {
		log.Printf("timeseries.shouldAppend: error reading last line: %s", err)
		return true
	}
	lastTimeStr := strings.Split(lastLine, ",")[0]
	lastTime, err := time.Parse(time.RFC3339, lastTimeStr)
	if err != nil {
		log.Printf("timeseries.shouldAppend: error parsing time of last line: %s", err)
		return true
	}

	return lastTime.Add(*onlyEvery).Before(clock.Now())
}
