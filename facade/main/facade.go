package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"net/http"
	"service"
)

func main() {
	panic(http.ListenAndServe(service.FacadeAddr, &FacadeListener{}))
}

type FacadeListener struct{}

func (m *FacadeListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responseBody string
	if request.Method == "GET" {
		responseBody = getLogs() + "\n" + getMessages()
	} else if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		var params service.RequestParams
		err = json.Unmarshal(body, &params)
		if err != nil {
			panic(err)
		}
		info := service.RequestInfo{Id: uuid.New().String(), Msg: params.Msg}
		logRequestMessage, _ := json.Marshal(info)
		err = sendToLogger(logRequestMessage)
		if err != nil {
			responseBody = err.Error()
		} else {
			responseBody = "Message sent successfully!"
		}
	}
	_, _ = writer.Write([]byte(responseBody))
}

func getLogs() string {
	for i := 0; i < service.LoggerPortsSize; i++ {
		logs, err := service.GetData(service.LoggingServiceAddr[i])
		if err != nil {
			fmt.Println("Logger on port " + service.LoggingServiceAddr[i] + " is not responding")
		} else {
			return logs
		}
	}
	return service.MsgServicesNotResponding
}

func getMessages() string {
	messages, _ := service.GetData(service.MessagesServiceAddr)
	return messages
}

func sendToLogger(message []byte) error {
	for tryCount := 0; tryCount < 15; tryCount++ {
		i := rand.Int() % 3
		_, err := http.Post(
			service.Localhost+service.LoggingServiceAddr[i],
			"application/json",
			bytes.NewReader(message))
		if err == nil {
			return nil
		} else {
			tryCount++
			fmt.Println("Cannot send message to logger " + service.LoggingServiceAddr[i])
		}
	}
	return errors.New(service.MsgServicesNotResponding)
}
