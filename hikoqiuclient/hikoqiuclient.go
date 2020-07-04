package hikoqiuclient

import (
	"eureka/vars"
	"fmt"
	"time"

	"github.com/HikoQiu/go-eureka-client/eureka"
)

// instance 정보에 meta 필드가 없음.

type HikoQiuClient struct {
	instanceID string
	vo         *eureka.InstanceVo
	config     *eureka.EurekaClientConfig
	cli        *eureka.Client
	api        *eureka.EurekaServerApi
	ticker     *time.Ticker
	done       chan bool
}

func NewClient() *HikoQiuClient {
	client := &HikoQiuClient{}

	client.vo = eureka.DefaultInstanceVo()
	client.vo.App = vars.AppName
	client.vo.Hostname = vars.Hostname
	client.vo.Status = eureka.STATUS_UP
	client.vo.Port.Value = vars.LocalPort
	client.vo.Port.Enabled = "true"

	return client
}

func (c *HikoQiuClient) Register() {
	c.config = eureka.GetDefaultEurekaClientConfig()
	c.config.UseDnsForFetchingServiceUrls = false
	c.config.AutoUpdateDnsServiceUrls = false
	c.config.Region = eureka.DEFAULT_REGION
	c.config.AvailabilityZones = map[string]string{
		eureka.DEFAULT_REGION: eureka.DEFAULT_ZONE,
	}
	c.config.ServiceUrl = map[string]string{
		eureka.DEFAULT_ZONE: vars.EurekaServerURL,
	}

	c.cli = eureka.DefaultClient.Config(c.config)
	c.api, _ = c.cli.Api()

	c.instanceID, _ = c.api.RegisterInstanceWithVo(c.vo)

	c.sendHeartBeat()
}

func (c *HikoQiuClient) sendHeartBeat() {
	c.ticker = time.NewTicker(time.Duration(vars.HeartBeatInterval) * time.Second)
	c.done = make(chan bool, 1)

	go func() {
		defer func() {
			fmt.Println("stop SendHeartBeat")
			c.ticker.Stop()
		}()

		for {
			select {
			case <-c.ticker.C:
				c.api.SendHeartbeat(vars.AppName, c.instanceID)
			case <-c.done:
				return
			}
		}
	}()
}

func (c *HikoQiuClient) DeRegister() {

	close(c.done)

	c.vo.Status = eureka.STATUS_DOWN
	c.instanceID, _ = c.api.RegisterInstanceWithVo(c.vo)

	c.api.DeRegisterInstance(vars.AppName, c.instanceID)
}
