package service

import (
	"io/ioutil"
	"net/http"
)

const Localhost = "http://127.0.0.1"
const FacadeAddr = ":31000"
const LoggingServiceAddr = ":31001"
const MessagesServiceAddr = ":31002"

type RequestInfo struct {
	Id  string
	Msg string
}

type RequestParams struct {
	Msg string
}

func GetData(serverAddress string) string {
	loggingResp, err := http.Get(Localhost + serverAddress)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(loggingResp.Body)
	return string(body)
}
