package common

import "errors"

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetConfigResponse struct {
	Message string  `json:"message"`
	Config  *Config `json:"configs"`
}

var ErrUpdateConfig = errors.New("update config failed")
