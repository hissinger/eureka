package fairwayclient

import (
	"eureka/vars"
	"strconv"

	"github.com/spectre013/fairway"
)

// port enable 필드가 bool로 정의 되어 있어서 에러남.

type FairwayClient struct {
	cli fairway.EurekaClient
}

func NewClient() *FairwayClient {
	return &FairwayClient{}
}

func (c *FairwayClient) Register() {
	conf := fairway.EurekaConfig{
		Name:        vars.AppName,
		URL:         vars.EurekaServerURL,
		VipAddress:  vars.AppName,
		IPAddress:   vars.LocalIP,
		HostName:    vars.LocalIP,
		Port:        strconv.Itoa(vars.LocalPort),
		SecurePort:  "443",
		RestService: false,
		PreferIP:    true,
	}

	c.cli = fairway.Init(conf)

	fairway.Register(conf) // Performs Eureka registration
}

func (c *FairwayClient) DeRegister() {
}
