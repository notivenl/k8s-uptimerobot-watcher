package main

import (
	"github.com/notivenl/uptime-kubernetes/internal/client"
	"github.com/notivenl/uptime-kubernetes/pkg/config"
	"github.com/notivenl/uptime-kubernetes/pkg/ports/kubernetes"
)

func main() {
	conf := config.NewConfig()

	c := client.NewUptimeClient(conf, kubernetes.NewRestClient())

	c.Start()
}
