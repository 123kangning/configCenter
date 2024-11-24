package main

import (
	"bytes"
	"cc/client/common"
	"cc/client/configs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestUpdateConfig(t *testing.T) {
	config := &common.Config{
		Key:   "k2",
		Value: "v3",
	}
	data, _ := json.Marshal(config)
	req, err := http.NewRequest("PUT", configs.ClientAddrPrefix+"/config", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	var answer map[string]interface{}
	err = json.Unmarshal(body, &answer)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if answer["message"] != "success" {
		t.Fatalf("Failed to update config: %v", answer)
	}
}

func TestGetConfig(t *testing.T) {

	req, err := http.NewRequest("GET", configs.ClientAddrPrefix+"/config", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	queryParams := fmt.Sprintf("key=%s&callback=%s", "k2", configs.ClientAddrPrefix+"/callback")
	req.URL.RawQuery = queryParams

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	getConfigResponse := &common.GetConfigResponse{}
	err = json.Unmarshal(body, getConfigResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	t.Log(string(body))
}
