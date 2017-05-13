package config

import (
	"fmt"
	"time"

	"github.com/jinzhu/configor"
)

type RedisConfig struct {
	Host string `default: localhost`
	Port uint   `default: 6379`
}

type QueueConfig struct {
	Pool        uint `default:"30"`
	Concurrency uint `default:"5"`
	Namespace   string
	Database    int
	Queues      []string
}

type JobConfig struct {
	Url string
}

type CronItemConfig struct {
	Name  string
	Url   string
	Times time.Duration `default: "5"`
}

var Config = struct {
	Redis RedisConfig
	Queue QueueConfig
	Job   JobConfig
	Crons []CronItemConfig
}{}

func init() {
	if err := configor.Load(&Config, "config.yml"); err != nil {
		fmt.Println("config load err: ", err)
	}
	fmt.Printf("config: %#v \n", Config)
}
