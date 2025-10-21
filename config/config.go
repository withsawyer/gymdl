package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfig 加载配置
func LoadConfig(file string) *Config {
	bytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("load config err:", err)
	}
	c := &Config{}
	err = json.Unmarshal(bytes, c)
	if err != nil {
		fmt.Println("parse config err:", err)
	}
	return c
}
