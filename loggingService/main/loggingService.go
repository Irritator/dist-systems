package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"service"
)

var messagesByIds = make(map[string]string)

func main() {
	panic(http.ListenAndServe(service.LoggingServiceAddr, &LoggingListener{}))
}

type LoggingListener struct{}

func (m *LoggingListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		logs, _ := json.Marshal(messagesByIds)
		_, _ = writer.Write(logs)
	} else {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		var requestInfo service.RequestInfo
		if err = json.Unmarshal(body, &requestInfo); err != nil {
			panic(err)
		} else {
			fmt.Println("requestInfo.Msg => " + requestInfo.Msg)
			fmt.Println("requestInfo.Id => " + requestInfo.Id)
		}
		messagesByIds[requestInfo.Id] = requestInfo.Msg
	}
}
