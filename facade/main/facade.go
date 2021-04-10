package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
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
		responseBody = service.GetData(service.LoggingServiceAddr) + "\n" + service.GetData(service.MessagesServiceAddr)
		//_, _ = writer.Write([]byte(responseBody))
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
		fmt.Println("sending message to logger")
		_, err = http.Post(
			service.Localhost+service.LoggingServiceAddr,
			"application/json",
			bytes.NewReader(logRequestMessage))
		if err != nil {
			responseBody = "There is an error processing your data. Check the input or contact your administrator"
		} else {
			responseBody = "Message sent successfully"
		}
	}
	_, _ = writer.Write([]byte(responseBody))
}
