package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
)

const MsgServicesNotResponding = "all logging services are not responding"
const QueueNameParam = "queueName"
const MessageQueueServiceName = "messageQueue"

const Localhost = "http://127.0.0.1"
const consulServices = Localhost + ":8500/v1/agent/service/"
const consulKeyValue = Localhost + ":8500/v1/kv/"

var Facades = []string{"facade1", "facade2"}
var Loggers = []string{"logger1", "logger2", "logger3"}
var Messengers = []string{"messenger1", "messenger2"}

type RequestInfo struct {
	Id  string
	Msg string
}

type RequestParams struct {
	Msg string
}

type ConsulServiceDTO struct {
	Id      string
	Service string
	Address string
	Port    int
	Meta    MetaDTO
}

type MetaDTO struct {
	HazelcastAddress string
}

type ValueDTO struct {
	Value string
}

func GetAvailablePort(serviceNames []string) string {
	for i := 0; i < len(serviceNames); i++ {
		port := GetServicePort(serviceNames[i])
		listener, err := net.Listen("tcp", port)
		if err != nil {
			fmt.Println("the port number", port, " is already taken, try next one")
		} else {
			_ = listener.Close()
			return port
		}
	}
	panic("no ports available")
}

func GetLoggerWithHazelcast() (string, string) {
	for i := 0; i < len(Loggers); i++ {
		serviceInfo := getServiceInfo(Loggers[i])
		port := ":" + strconv.Itoa(serviceInfo.Port)
		listener, err := net.Listen("tcp", port)
		if err != nil {
			fmt.Println("the port number", port, " is already taken, try next one")
		} else {
			_ = listener.Close()
			return port, serviceInfo.Meta.HazelcastAddress
		}
	}
	panic("no ports available")
}

func GetRandomAddress(serviceNames []string) string {
	index := rand.Int() % len(serviceNames)
	return GetServiceAddress(serviceNames[index])
}

func GetServicePort(serviceName string) string {
	serviceInfo := getServiceInfo(serviceName)
	return ":" + strconv.Itoa(serviceInfo.Port)
}

func GetServiceAddress(serviceName string) string {
	serviceInfo := getServiceInfo(serviceName)
	if serviceInfo.Address == "" {
		serviceInfo.Address = Localhost
	}
	return serviceInfo.Address + ":" + strconv.Itoa(serviceInfo.Port)
}

func getServiceInfo(serviceName string) *ConsulServiceDTO {
	resp, _ := http.Get(consulServices + serviceName)
	if resp.StatusCode == 404 {
		panic("no addresses available")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var serviceInfo ConsulServiceDTO
	_ = json.Unmarshal(body, &serviceInfo)
	return &serviceInfo
}

func GetConsulValue(key string) string {
	resp, _ := http.Get(consulKeyValue + key)
	body, _ := ioutil.ReadAll(resp.Body)
	var values []ValueDTO
	_ = json.Unmarshal(body, &values)
	decodedMsg, _ := base64.StdEncoding.DecodeString(values[0].Value)
	fmt.Println("key ===> ", key)
	fmt.Println("value ===> ", string(decodedMsg))
	return string(decodedMsg)
}

func Get(serverAddress string) (string, error) {
	loggingResp, err := http.Get(serverAddress)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(loggingResp.Body)
	return string(body), nil
}
