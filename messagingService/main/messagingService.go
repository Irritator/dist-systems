package main

import (
	"net/http"
	"service"
)

func main() {
	panic(http.ListenAndServe(service.MessagesServiceAddr, &MessagingListener{}))
}

type MessagingListener struct{}

func (m *MessagingListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		_, _ = writer.Write([]byte("Messaging Service currently unavailable"))
	}
	//todo currently unavailable
}
