package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func readFile(filePath string) (string, error) {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
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
