//  go test -v -bench=. benchmark_test.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func httpLogin() {
	resp, err := http.Post("http://127.0.0.1:8080/login",
		"application/x-www-form-urlencoded",
		strings.NewReader("{\"name\":\"tom\", \"password\":\"123456\"}"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func httpChat() {
	resp, err := http.Post("http://127.0.0.1:8080/chat",
		"application/x-www-form-urlencoded",
		strings.NewReader("{\"user_id\":\"3ees4489\", \"msg\":\"hello world!\"}"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func Benchmark_HttpServer(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n++
		if i%2 == 0 {
			httpLogin()
		} else {
			httpChat()
		}
	}
}
