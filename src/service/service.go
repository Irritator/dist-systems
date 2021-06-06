package service

import (
	"io/ioutil"
	"net/http"
)

const MsgServicesNotResponding = "all logging services are not responding"

const Localhost = "http://127.0.0.1"
const FacadeAddr = ":31000"

var LoggingServiceAddr = [...]string{":31101", ":31102", ":31103"}
var HazelcastAddr = [...]string{":5701", ":5702", ":5703"}
var MessagesServiceAddr = [...]string{":31002", ":31003"}

var LoggerPortsSize = len(LoggingServiceAddr)

type RequestInfo struct {
	Id  string
	Msg string
}

type RequestParams struct {
	Msg string
}

func Get(serverAddress string) (string, error) {
	loggingResp, err := http.Get(Localhost + serverAddress)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(loggingResp.Body)
	return string(body), nil
}
