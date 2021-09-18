package config

import "os"

type List struct {
	ConsulHost string `json:"consul_host"`
	ConsulPort string `json:"consul_port"`
}

func (list *List) Get() {
	list.ConsulHost, _ = os.LookupEnv("CONSUL_HOST")
	list.ConsulPort, _ = os.LookupEnv("CONSUL_PORT")
}
