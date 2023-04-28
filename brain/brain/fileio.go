package brain

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func readFile(filePath string) (string, error) {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading file %s: %s\n", filePath, err)
		return "", err
	}

	return string(buffer), err
}

func writeToFile(filePath string, value string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.WriteString(value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()
}

func appendToFile(filePath string, text string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s for append: %s", filePath, err)
	}
	defer file.Close()

	if _, err = file.WriteString("\n" + text); err != nil {
		return fmt.Errorf("failed to write to file %s: %s", filePath, err)
	}
	return nil
}

func getLastLineWithSeek(filePath string) (string, error) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, err := fileHandle.Stat()
	if err != nil {
		return "", err
	}

	filesize := stat.Size()
	for {
		cursor -= 1
		_, err := fileHandle.Seek(cursor, io.SeekEnd)
		if err != nil {
			return "", err
		}

		char := make([]byte, 1)
		_, err = fileHandle.Read(char)
		if err != nil {
			return "", err
		}

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}

	return line, nil
}
