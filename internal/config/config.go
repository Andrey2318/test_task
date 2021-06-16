package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	GRPCAddr  string   `json:"grpc_addr"`
	ProxyAddr string   `json:"proxy_addr"`
	ProxyPool []string `json:"proxy_pool"`
}

func New(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(file, c); err != nil {
		return nil, err
	}

	return c, nil
}
