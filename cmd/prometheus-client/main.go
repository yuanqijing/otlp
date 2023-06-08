package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Replace with your application's /metrics endpoint
	url := "http://localhost:8080/metrics"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to send request to %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v\n", err)
		return
	}

	fmt.Printf("Response:\n%s\n", body)
}
