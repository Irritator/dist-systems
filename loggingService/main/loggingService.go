package main

import (
	"encoding/json"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/core"
	"io/ioutil"
	"net/http"
	"service"
)

var logs core.Map

func main() {
	serviceInfo := service.RegisterService(service.Logger)
	logs = getHazelcastMap(serviceInfo.Meta.HazelcastAddress)
	panic(http.ListenAndServe(serviceInfo.GetStringPort(), &LoggingListener{}))
}

type LoggingListener struct{}

func (m *LoggingListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		valuePairs, _ := logs.EntrySet()
		msg := ""
		if len(valuePairs) > 0 {
			for i := 0; i < len(valuePairs); i++ {
				msg += valuePairs[i].Key().(string) + " : " + valuePairs[i].Value().(string) + "\n"
			}
		} else {
			msg = "There are no log messages"
		}
		_, _ = writer.Write([]byte(msg))
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
		_, _ = logs.Put(requestInfo.Id, requestInfo.Msg)
	}
}

func getHazelcastMap(address string) core.Map {
	config := hazelcast.NewConfig()
	config.NetworkConfig().AddAddress(address)
	client, err := hazelcast.NewClientWithConfig(config)
	if err != nil {
		fmt.Println(err)
	}
	logMap, _ := client.GetMap("log")
	return logMap
}
