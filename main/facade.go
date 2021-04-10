package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

type RequestInfo struct {
	Id  string
	Msg string
}

type RequestParams struct {
	Msg string
}

type FacadeListener struct{}

func (m *FacadeListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("==========================")
	fmt.Println("FacadeListener")
	if request.Method == "GET" {
		responseBody := getData(loggingServiceAddr) + "\n" + getData(messagesServiceAddr)
		_, _ = writer.Write([]byte(responseBody))
	} else if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		var params RequestParams
		err = json.Unmarshal(body, &params)
		if err != nil {
			panic(err)
		}
		info := RequestInfo{uuid.New().String(), params.Msg}
		logRequestMessage, _ := json.Marshal(info)
		_, _ = http.Post(
			"http://127.0.0.1"+loggingServiceAddr,
			"application/json",
			bytes.NewReader(logRequestMessage))
	}
}

func getData(serverAddress string) string {
	loggingResp, err := http.Get(serverAddress)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(loggingResp.Body)
	return string(body)
}
