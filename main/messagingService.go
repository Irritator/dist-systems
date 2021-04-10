package main

import (
	"fmt"
	"net/http"
)

type MessagingListener struct{}

func (m *MessagingListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("==========================")
	fmt.Println("MessagingListener")
	//todo currently unavailable
}
