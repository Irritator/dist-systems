package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"service"
)

var ch *amqp.Channel
var messageQueue amqp.Queue

func main() {
	queueName := service.GetConsulValue(service.QueueNameParam)
	conn, _ := amqp.Dial(service.GetServiceAddress(service.MessageQueueServiceName))
	ch, _ = conn.Channel()
	messageQueue, _ = ch.QueueDeclare(queueName, true, false, false, false, nil)
	port := service.GetAvailablePort(service.Facades)
	panic(http.ListenAndServe(port, &FacadeListener{}))
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
	for i := 0; i < len(service.Loggers); i++ {
		address := service.GetServiceAddress(service.Loggers[i])
		logs, err := service.Get(address)
		if err != nil {
			fmt.Println("Logger on port " + address + " is not responding")
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
	messageServiceAddress := service.GetRandomAddress(service.Messengers)
	messages, _ := service.Get(messageServiceAddress)
	return messages
}

func sendToLogger(reqParams service.RequestParams) error {
	info := service.RequestInfo{Id: uuid.New().String(), Msg: reqParams.Msg}
	logRequestMessage, _ := json.Marshal(info)
	for tryCount := 0; tryCount < 15; tryCount++ {
		loggerAddress := service.GetRandomAddress(service.Loggers)
		fmt.Println("loggerAddress ===> ", loggerAddress)
		_, err := http.Post(
			loggerAddress,
			"application/json",
			bytes.NewReader(logRequestMessage))
		if err == nil {
			return nil
		} else {
			tryCount++
			fmt.Println("Cannot send message to logger " + loggerAddress)
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
