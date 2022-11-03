package config

import (
	"os"
	"sync"
)

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var instance *StorageConfig
var once sync.Once

func GetConfig() *StorageConfig {
	once.Do(func() {
		instance = &StorageConfig{
			Host:     os.Getenv("host"),
			Port:     os.Getenv("port"),
			Database: os.Getenv("database"),
			Username: os.Getenv("user"),
			Password: os.Getenv("password")}

	})
	return instance
}
