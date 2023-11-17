package fileio

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ReadFile(filePath string) (string, error) {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading file %s: %s\n", filePath, err)
		return "", err
	}

	return string(buffer), err
}

func ReadLastLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}

func WriteToFile(filePath string, value string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	n, err := f.WriteString(value)
	if err != nil {
		return err
	}

	fmt.Printf("wrote %d bytes\n", n)
	err = f.Sync()
	return err
}

func AppendLineToFile(filePath, line string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s for append: %s", filePath, err)
	}
	defer file.Close()

	if _, err = file.WriteString("\n" + line); err != nil {
		return fmt.Errorf("failed to write to file %s: %s", filePath, err)
	}
	return nil
}
