package springbootclient

import (
	"eureka/vars"
	"fmt"

	"github.com/pineda89/golang-springboot/config"
	"github.com/pineda89/golang-springboot/eureka"
)

type SprintBootClient struct {
}

func NewClient() *SprintBootClient {
	return &SprintBootClient{}
}

func (c *SprintBootClient) Register() {
	config.LoadConfig()

	fmt.Println(config.Configuration)

	config.Configuration["eureka.client.serviceUrl.defaultZone"] = vars.EurekaServerURL
	config.Configuration["eureka.instance.ip-address"] = vars.LocalIP
	config.Configuration["spring.application.name"] = vars.AppName
	config.Configuration["server.port"] = vars.LocalPort
	config.Configuration["hostname"] = vars.LocalIP

	eureka.Register(config.Configuration)
}

func (c *SprintBootClient) DeRegister() {
	eureka.Deregister()
}
