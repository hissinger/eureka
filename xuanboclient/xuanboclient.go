package xuanboclient

import (
	"eureka/vars"
	"fmt"
	"time"

	eureka "github.com/xuanbo/eureka-client"
)

type XuanboClient struct {
	client   *eureka.Client
	instance *eureka.Instance
	ticker   *time.Ticker
	done     chan bool
}

func NewClient() *XuanboClient {
	return &XuanboClient{}
}

func (c *XuanboClient) Register() {

	conf := &eureka.Config{
		DefaultZone:           vars.EurekaServerURL + "/",
		App:                   vars.AppName,
		Port:                  vars.LocalPort,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"management.port": "8080",
		},
	}

	// create eureka client
	c.client = eureka.NewClient(conf)
	// start client, register、heartbeat、refresh
	// c.client.Start()

	c.instance = eureka.NewInstance(vars.LocalIP, conf)
	eureka.Register(c.client.Config.DefaultZone, c.client.Config.App, c.instance)

	c.sendHeartBeat()
}

func (c *XuanboClient) sendHeartBeat() {
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
				eureka.Heartbeat(c.client.Config.DefaultZone, vars.AppName, c.instance.InstanceID)
			case <-c.done:
				return
			}
		}
	}()
}

func (c *XuanboClient) DeRegister() {
	close(c.done)

	c.instance.Status = "DOWN"
	eureka.Register(c.client.Config.DefaultZone, vars.AppName, c.instance)
	eureka.UnRegister(c.client.Config.DefaultZone, vars.AppName, c.instance.InstanceID)
}
