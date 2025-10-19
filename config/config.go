package config

import (
	"encoding/json"
	"os"

	"github.com/nichuanfang/gymdl/utils"
)

// LoadConfig 加载配置
func LoadConfig(file string) *Config {
	bytes, err := os.ReadFile(file)
	if err != nil {
		utils.Logger.FatalF("load config err: %v", err)
	}

	c := &Config{}
	err = json.Unmarshal(bytes, c)
	if err != nil {
		utils.Logger.FatalF("parse config err: %v", err)
	}

	return c
}
