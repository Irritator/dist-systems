package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"net/http"
	"service"
	"strings"
)

var ch *amqp.Channel
var messageQueue amqp.Queue
var messages []string

func main() {
	queueName := service.GetConsulValue(service.QueueNameParam)
	conn, _ := amqp.Dial(service.GetServiceAddress(service.MessageQueueServiceName))
	ch, _ = conn.Channel()
	messageQueue, _ = ch.QueueDeclare(queueName, true, false, false, false, nil)
	msgs, _ := ch.Consume(messageQueue.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			msg := string(d.Body)
			messages = append(messages, msg)
			fmt.Println("Received a message:", msg)
		}
	}()

	port := service.GetAvailablePort(service.Messengers)
	panic(http.ListenAndServe(port, &MessagingListener{}))
}

type MessagingListener struct{}

func (m *MessagingListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		text := strings.Join(messages, "\n")
		_, _ = writer.Write([]byte(text))
	}

}
