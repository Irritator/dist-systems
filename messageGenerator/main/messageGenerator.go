package main

import (
	"bytes"
	"net/http"
	"service"
	"strconv"
)

func main() {
	for i := 0; i < 5; i++ {
		message := `{"msg": "test message` + strconv.Itoa(i) + `"}`
		_, _ = http.Post(
			service.Localhost+":40000",
			"application/json",
			bytes.NewReader([]byte(message)))
	}
	for i := 5; i < 10; i++ {
		message := `{"msg": "test message` + strconv.Itoa(i) + `"}`
		_, _ = http.Post(
			service.Localhost+":40001",
			"application/json",
			bytes.NewReader([]byte(message)))
	}
}
