package registry

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

type client struct {
	client *api.Client
}

func NewClient() *client {
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("api.NewClient failed, err:%v\n", err)
	}
	return &client{client: consulClient}
}

func (c *client) RegisterService(serviceName, ip string, port int) error {
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ip, port),
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "1m",
	}
	srv := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
		Name:    serviceName,
		Tags:    []string{serviceName},
		Address: ip,
		Port:    port,
		Check:   check,
	}
	return c.client.Agent().ServiceRegister(srv)
}

func (c *client) ListService(serviceName string) (map[string]*api.AgentService, error) {
	return c.client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, serviceName))
}

func (c *client) DeregisterService(serviceId string) error {
	return c.client.Agent().ServiceDeregister(serviceId)
}
