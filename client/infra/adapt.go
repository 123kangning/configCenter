package infra

import (
	"bytes"
	"cc/client/common"
	"cc/client/configs"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func updateConfig(config *common.Config) error {
	requestData, _ := json.Marshal(config)
	req, err := http.NewRequest("PUT", configs.ServerAddrPrefix+"/config", bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var answer map[string]interface{}
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return err
	}
	if answer["message"] != "success" {
		return common.ErrUpdateConfig
	}
	return nil
}

func getRemoteConfig(key string) (*common.Config, error) {
	queryParams := fmt.Sprintf("key=%s&callback=%s", key, configs.ClientAddrPrefix+"/callback")
	request, err := http.NewRequest("GET", configs.ServerAddrPrefix+"/config", nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = queryParams
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	getConfigResponse := &common.GetConfigResponse{}
	err = json.Unmarshal(body, getConfigResponse)
	if err != nil {
		return nil, err
	}
	if getConfigResponse.Message != "success" {
		return nil, errors.New(getConfigResponse.Message)
	}
	return getConfigResponse.Config, nil
}
