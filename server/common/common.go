package common

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetConfigResponse struct {
	Message string  `json:"message"`
	Config  *Config `json:"configs"`
}
