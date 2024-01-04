package config

import (
	"os"
)

type Config struct {
	Db struct {
		Mongo struct {
			Url string `json:"url"`
		} `json:"mongo"`
		Yugabyte struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Dbname   string `json:"dbname"`
			Address  string `json:"address"`
			PoolSize string `json:"pool_size"`
			Port     string `json:"port"`
		} `json:"yugabyte"`
	} `json:"db"`
}

func GetConfig() *Config {
	var config Config

	if yugabyteURI := os.Getenv("YUGABYTE_URI"); yugabyteURI != "" {
		config.Db.Yugabyte.Address = yugabyteURI
		config.Db.Yugabyte.Password = "yugabyte"
		if os.Getenv("YUGABYTE_PASSWORD") != "" {
			config.Db.Yugabyte.Password = os.Getenv("YUGABYTE_PASSWORD")
		}
		config.Db.Yugabyte.Username = "yugabyte"
		if os.Getenv("YUGABYTE_USERNAME") != "" {
			config.Db.Yugabyte.Username = os.Getenv("YUGABYTE_USERNAME")
		}
		config.Db.Yugabyte.Dbname = "yugabyte"
		if os.Getenv("YUGABYTE_DBNAME") != "" {
			config.Db.Yugabyte.Dbname = os.Getenv("YUGABYTE_DBNAME")
		}
		config.Db.Yugabyte.PoolSize = "20"
		if os.Getenv("POOL_SIZE") != "" {
			config.Db.Yugabyte.PoolSize = os.Getenv("POOL_SIZE")
		}
	}

	if mongoURI := os.Getenv("MONGODB_URI"); mongoURI != "" {
		config.Db.Mongo.Url = mongoURI
	}
	return &config
}
