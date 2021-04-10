package main

import "net/http"

const facadeAddr = ":30000"
const loggingServiceAddr = ":30001"
const messagesServiceAddr = ":30002"

func main() {
	go func() {
		panic(http.ListenAndServe(loggingServiceAddr, &LoggingListener{}))
	}()
	go func() {
		panic(http.ListenAndServe(messagesServiceAddr, &MessagingListener{}))
	}()
	panic(http.ListenAndServe(facadeAddr, &FacadeListener{}))
}
