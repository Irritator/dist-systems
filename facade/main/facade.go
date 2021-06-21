package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"service"
)

var ch *amqp.Channel
var messageQueue amqp.Queue

func main() {
	serviceInfo := service.RegisterService(service.Facade)
	queueName := service.GetConsulValue(service.QueueNameParam)
	conn, _ := amqp.Dial(service.GetConsulValue(service.MessageQueueServiceAddress))
	ch, _ = conn.Channel()
	messageQueue, _ = ch.QueueDeclare(queueName, true, false, false, false, nil)
	panic(http.ListenAndServe(serviceInfo.GetStringPort(), &FacadeListener{}))
}

type FacadeListener struct{}

func (m *FacadeListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var responseBody string
	if request.Method == "GET" {
		responseBody = getLogs() + "\n" + getMessages()
	} else if request.Method == "POST" {
		reqParams := parseRequest(request)
		sendToLogger(reqParams)
		sendToMessenger(reqParams)
	}
	_, _ = writer.Write([]byte(responseBody))
}

func getLogs() string {
	loggerService := service.GetRandomService(service.Logger)
	logs, _ := service.Get(loggerService.GetFullAddress())
	return logs
}

func parseRequest(request *http.Request) service.RequestParams {
	body, _ := ioutil.ReadAll(request.Body)
	var params service.RequestParams
	_ = json.Unmarshal(body, &params)
	return params
}

func getMessages() string {
	messenger := service.GetRandomService(service.Messenger)
	messages, _ := service.Get(messenger.GetFullAddress())
	return messages
}

func sendToLogger(reqParams service.RequestParams) {
	info := service.RequestInfo{Id: uuid.New().String(), Msg: reqParams.Msg}
	logRequestMessage, _ := json.Marshal(info)
	loggerService := service.GetRandomService(service.Logger)
	_, _ = http.Post(
		loggerService.GetFullAddress(),
		"application/json",
		bytes.NewReader(logRequestMessage))
}

func sendToMessenger(reqParams service.RequestParams) {
	_ = ch.Publish(
		"", messageQueue.Name, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(reqParams.Msg),
		})
}
