package main

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {

	req := httptest.NewRequest("GET", "/", nil)
	makeReq := httptest.NewRecorder()
	handler(makeReq, req)

	response := makeReq.Result()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	if response.StatusCode == 200 {
		t.Log("Received HTTP 200 OK.")
	} else {
		t.Error("Test request failed.")
	}
}
