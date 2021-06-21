package service

import (
	"bytes"
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
const Localhost = "http://127.0.0.1"
const consulServiceDescribe = Localhost + ":8500/v1/catalog/service/"
const consulServiceRegister = Localhost + ":8500/v1/agent/service/register"
const consulKeyValue = Localhost + ":8500/v1/kv/"
const hazelcast = "hazelcast/"

const Facade = "facade"
const Logger = "logger"
const Messenger = "messenger"
const MessageQueueServiceAddress = "messageQueueAddress"
const QueueNameParam = "queueName"

type RequestInfo struct {
	Id  string
	Msg string
}

type RequestParams struct {
	Msg string
}

type ConsulServiceDTO struct {
	Id      string
	Name    string
	Address string
	Port    int
	Meta    MetaDTO
}

type ConsulServiceGetDTO struct {
	ServiceID      string
	ServiceName    string
	ServiceAddress string
	ServicePort    int
	ServiceMeta    MetaDTO
}

func (serviceDto ConsulServiceDTO) GetStringPort() string {
	return ":" + strconv.Itoa(serviceDto.Port)
}
func (serviceDto ConsulServiceDTO) GetFullAddress() string {
	return serviceDto.Address + serviceDto.GetStringPort()
}

func (serviceDto ConsulServiceGetDTO) GetStringPort() string {
	return ":" + strconv.Itoa(serviceDto.ServicePort)
}
func (serviceDto ConsulServiceGetDTO) GetFullAddress() string {
	return serviceDto.ServiceAddress + serviceDto.GetStringPort()
}

type MetaDTO struct {
	HazelcastAddress string
}

type ValueDTO struct {
	Value string
}

func RegisterService(serviceName string) ConsulServiceDTO {
	port := GetAvailablePort()
	existingServices := getServices(serviceName)
	index := len(existingServices) + 1
	serviceId := serviceName + strconv.Itoa(index)
	var meta MetaDTO
	if serviceName == Logger {
		meta.HazelcastAddress = GetConsulValue(hazelcast + strconv.Itoa(index))
	}
	return registerService(serviceId, serviceName, Localhost, port, meta)
}

func registerService(id string, serviceName string, address string, port int, meta MetaDTO) ConsulServiceDTO {
	serviceDescribe := ConsulServiceDTO{
		Id:      id,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Meta:    meta,
	}
	bodyJson, _ := json.Marshal(serviceDescribe)
	req, _ := http.NewRequest(http.MethodPut, consulServiceRegister, bytes.NewBuffer(bodyJson))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	res, _ := client.Do(req)
	fmt.Println(string(bodyJson))
	fmt.Println(res)
	fmt.Println("Starting service " + serviceName + " on port " + serviceDescribe.GetStringPort())
	return serviceDescribe
}

func getServices(serviceName string) []ConsulServiceDTO {
	resp, _ := http.Get(consulServiceDescribe + serviceName)
	body, _ := ioutil.ReadAll(resp.Body)
	var serviceInfos []ConsulServiceGetDTO
	_ = json.Unmarshal(body, &serviceInfos)
	var convertedServiceInfos []ConsulServiceDTO
	for i := 0; i < len(serviceInfos); i++ {
		convertedServiceInfos = append(convertedServiceInfos, ConsulServiceDTO{
			Id:      serviceInfos[i].ServiceID,
			Name:    serviceInfos[i].ServiceName,
			Address: serviceInfos[i].ServiceAddress,
			Port:    serviceInfos[i].ServicePort,
			Meta:    serviceInfos[i].ServiceMeta,
		})
	}
	return convertedServiceInfos
}

func GetRandomService(serviceName string) ConsulServiceDTO {
	services := getServices(serviceName)
	if len(services) == 0 {
		panic("No services available: " + serviceName)
	}
	return services[rand.Int()%len(services)]
}

func GetAvailablePort() int {
	for port := 40000; port < 65535; port++ {
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			fmt.Println("the port number", port, " is already taken, try next one")
		} else {
			_ = listener.Close()
			return port
		}
	}
	panic("no ports available")
}

func GetConsulValue(key string) string {
	fmt.Println("get value from Cosul: " + key)
	resp, _ := http.Get(consulKeyValue + key)
	body, _ := ioutil.ReadAll(resp.Body)
	var values []ValueDTO
	_ = json.Unmarshal(body, &values)
	if len(values) == 0 {
		panic("value is missing in consul: " + key)
	}
	decodedMsg, _ := base64.StdEncoding.DecodeString(values[0].Value)
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
