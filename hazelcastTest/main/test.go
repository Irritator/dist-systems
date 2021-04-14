package main

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/core"
	"time"
)

func main() {
	//initMap()
	//multipleConnectionsWithoutLock()
	//multipleConnectionsPessimisticLock()
	//multipleConnectionsOptimisticLock()
	processQueue()
}

func multipleConnectionsWithoutLock() {
	go updateWithoutLock("5701")
	go updateWithoutLock("5702")
	updateWithoutLock("5703")
}

func updateWithoutLock(port string) {
	testMap := getMap(port)
	fmt.Println("map got on port " + port)
	for i := 0; i < 1000; i++ {
		value, _ := testMap.Get(i)
		newVal := value.(int64) + 1
		time.Sleep(10 * time.Millisecond)
		oldVal, _ := testMap.Put(i, newVal)
		fmt.Println(i, " update SUCCESS on port "+port, " value is: ", newVal, "old value is: ", oldVal, "//old must be ", value)
	}
}

func multipleConnectionsPessimisticLock() {
	go updateWithPessimisticLock("5701")
	go updateWithPessimisticLock("5702")
	updateWithPessimisticLock("5703")
}

func updateWithPessimisticLock(port string) {
	testMap := getMap(port)
	fmt.Println("map got on port " + port)
	for i := 0; i < 1000; i++ {
		_ = testMap.Lock(i)
		fmt.Println(i, " locked on port "+port)
		value, _ := testMap.Get(i)
		fmt.Println(i, "oldvalue ", value, " retrieved on port "+port)
		newVal := value.(int64) + 1
		time.Sleep(20 * time.Millisecond)
		oldVal, _ := testMap.Put(i, newVal)
		err := testMap.Unlock(i)
		if err != nil {
			fmt.Println(i, " cannot write on port "+port)
		} else {
			fmt.Println(i, " update SUCCESS on port "+port, " value is: ", newVal, "old value is: ", oldVal, "//old must be ", value)
		}
	}
}

func multipleConnectionsOptimisticLock() {
	go updateWithOptimisticLock("5701")
	go updateWithOptimisticLock("5702")
	updateWithOptimisticLock("5703")
}

func updateWithOptimisticLock(port string) {
	testMap := getMap(port)
	fmt.Println("map got on port " + port)
	for i := 0; i < 1000; i++ {
		for {
			value, _ := testMap.Get(i)
			fmt.Println(i, " || oldvalue ", value, " retrieved on port "+port)
			newVal := value.(int64) + 1
			time.Sleep(20 * time.Millisecond)
			isReplaced, _ := testMap.ReplaceIfSame(i, value, newVal)
			if isReplaced {
				fmt.Println(i, " || port: "+port, " || value updated successfully from ", value, " to ", newVal)
				break
			} else {
				fmt.Println(i, " || port: "+port, " || value changed during transaction. Try again")
			}
		}
	}
}

func initMap() {
	client, _ := hazelcast.NewClient()
	hazelcastMap, _ := client.GetMap("map")
	for i := 0; i < 1000; i++ {
		hazelcastMap.Put(i, i)

	}
	client.Shutdown()
}

func processQueue() {
	//go readFromQueue("5702")
	//go readFromQueue("5703")
	writeToQueue("5701")
}

func writeToQueue(port string) {
	testQueue := getQueue(port)
	for i := 0; i < 100; i++ {
		err := testQueue.Put(i)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(i, "added")
		time.Sleep(100)
	}
}

func readFromQueue(port string) {
	testQueue := getQueue(port)
	for {
		index, _ := testQueue.Take()
		fmt.Println(index, "is consumed on port "+port)
		time.Sleep(100)
	}
}

func getMap(port string) core.Map {
	client := getClient(port)
	testMap, _ := client.GetMap("map")
	return testMap
}

func getQueue(port string) core.Queue {
	client := getClient(port)
	testQueue, _ := client.GetQueue("testQueue")
	return testQueue

}

func getClient(port string) hazelcast.Client {
	config := hazelcast.NewConfig()
	config.NetworkConfig().AddAddress("127.0.0.1:" + port)
	client, _ := hazelcast.NewClientWithConfig(config)
	return client
}
