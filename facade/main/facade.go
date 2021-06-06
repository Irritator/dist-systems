package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"io/ioutil"
	"math/rand"
	"net/http"
	"service"
)

var ch *amqp.Channel
var messageQueue amqp.Queue

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ = conn.Channel()
	messageQueue, _ = ch.QueueDeclare("messaging_service", true, false, false, false, nil)
	panic(http.ListenAndServe(service.FacadeAddr, &FacadeListener{}))
}

type FacadeListener struct{}

func (m *FacadeListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responseBody string
	if request.Method == "GET" {
		responseBody = getLogs() + "\n" + getMessages()
	} else if request.Method == "POST" {
		reqParams := parseRequest(request)
		err := sendToLogger(reqParams)
		sendToMessenger(reqParams)
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
		logs, err := service.Get(service.LoggingServiceAddr[i])
		if err != nil {
			fmt.Println("Logger on port " + service.LoggingServiceAddr[i] + " is not responding")
		} else {
			return logs
		}
	}
	return service.MsgServicesNotResponding
}

func parseRequest(request *http.Request) service.RequestParams {
	body, _ := ioutil.ReadAll(request.Body)
	var params service.RequestParams
	_ = json.Unmarshal(body, &params)
	return params
}

func getMessages() string {
	i := rand.Int() % 3
	messageServiceAddress := service.MessagesServiceAddr[i]
	messages, _ := service.Get(messageServiceAddress)
	return messages
}

func sendToLogger(reqParams service.RequestParams) error {
	info := service.RequestInfo{Id: uuid.New().String(), Msg: reqParams.Msg}
	logRequestMessage, _ := json.Marshal(info)
	for tryCount := 0; tryCount < 15; tryCount++ {
		i := rand.Int() % 3
		_, err := http.Post(
			service.Localhost+service.LoggingServiceAddr[i],
			"application/json",
			bytes.NewReader(logRequestMessage))
		if err == nil {
			return nil
		} else {
			tryCount++
			fmt.Println("Cannot send message to logger " + service.LoggingServiceAddr[i])
		}
	}
	return errors.New(service.MsgServicesNotResponding)
}

func sendToMessenger(reqParams service.RequestParams) {
	_ = ch.Publish(
		"", messageQueue.Name, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(reqParams.Msg),
		})
}
