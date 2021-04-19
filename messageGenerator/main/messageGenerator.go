package main

import (
	"bytes"
	"net/http"
	"service"
	"strconv"
)

func main() {
	for i := 0; i < 10; i++ {
		message := `{"msg": "test message` + strconv.Itoa(i) + `"}`
		_, _ = http.Post(
			service.Localhost+service.FacadeAddr,
			"application/json",
			bytes.NewReader([]byte(message)))
	}
}
