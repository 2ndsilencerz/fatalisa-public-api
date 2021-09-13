package config

import "os"

type List struct {
	RedisHost string `json:"redis_host"`
	RedisPass string `json:"redis_pass"`

	ConsulHost string `json:"consul_host"`
	ConsulPort string `json:"consul_port"`
}

func (list *List) Get() {
	list.RedisHost, _ = os.LookupEnv("REDIS_HOST")
	list.RedisPass, _ = os.LookupEnv("REDIS_PASS")

	list.ConsulHost, _ = os.LookupEnv("CONSUL_HOST")
	list.ConsulPort, _ = os.LookupEnv("CONSUL_PORT")
}
