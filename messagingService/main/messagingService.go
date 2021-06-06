package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"net"
	"net/http"
	"service"
	"strings"
)

var ch *amqp.Channel
var messageQueue amqp.Queue
var messages []string

func main() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ = conn.Channel()
	messageQueue, _ = ch.QueueDeclare("messaging_service", true, false, false, false, nil)
	msgs, _ := ch.Consume(messageQueue.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			msg := string(d.Body)
			messages = append(messages, msg)
			fmt.Println("Received a message:", msg)
		}
	}()

	panic(http.ListenAndServe(findAvailablePort(), &MessagingListener{}))
}

type MessagingListener struct{}

func (m *MessagingListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		text := strings.Join(messages, "\n")
		_, _ = writer.Write([]byte(text))
	}

}

func findAvailablePort() string {
	for i := 0; i < len(service.MessagesServiceAddr); i++ {
		listener, err := net.Listen("tcp", service.MessagesServiceAddr[i])
		if err != nil {
			fmt.Println(service.MessagesServiceAddr[i], " is already taken, try next one")
		} else {
			_ = listener.Close()
			return service.MessagesServiceAddr[i]
		}
	}
	panic("All ports are unavailable")
}
