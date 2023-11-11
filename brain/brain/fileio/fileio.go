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

func WriteToFile(filePath string, value string) {
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

func AppendLineToFile(filePath, line string) error {
	// Open the file in append mode, creating it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("\n" + line)
	if err != nil {
		return err
	}

	return nil
}
