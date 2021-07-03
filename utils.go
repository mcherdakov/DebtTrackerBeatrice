package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func formRequestBuffer(requestBody map[string]interface{}) *bytes.Buffer {
	bytesRepresentation, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalln(err)
	}

	return bytes.NewBuffer(bytesRepresentation)
}

func postRequest(handler string, requestBody map[string]interface{}) []byte {
	resp, err := http.Post(
		fmt.Sprintf(
			"%s/%s",
			telegramUrl,
			handler,
		),
		"application/json",
		formRequestBuffer(requestBody),
	)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}
