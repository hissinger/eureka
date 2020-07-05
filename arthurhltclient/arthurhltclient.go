package arthurhltclient

import (
	"eureka/vars"
	"strconv"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

// port enable 필드가 bool로 정의 되어 있어서 에러남.

type ArthurHltClient struct {
	client   *eureka.Client
	instance *eureka.InstanceInfo
	ticker   *time.Ticker
}

func NewClient() *ArthurHltClient {
	return &ArthurHltClient{}
}

func (c *ArthurHltClient) Register() {
	c.client = eureka.NewClient([]string{vars.EurekaServerURL})
	c.instance = eureka.NewInstanceInfo(vars.LocalIP, vars.AppName, vars.LocalIP, vars.LocalPort, uint(vars.HeartBeatInterval), false) //Create a new instance to register
	c.instance.InstanceID = vars.LocalIP + ":" + vars.AppName + ":" + strconv.Itoa(vars.LocalPort)
	c.instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	c.instance.Metadata.Map["foo"] = "bar"              //add metadata for example
	c.client.RegisterInstance(vars.AppName, c.instance) // Register new instance in your eureka(s)

	go c.sendHeartbeat()
}

func (c *ArthurHltClient) sendHeartbeat() {
	c.ticker = time.NewTicker(time.Duration(vars.HeartBeatInterval) * time.Second)
	for {
		select {
		case <-c.ticker.C:
			c.client.SendHeartbeat(c.instance.App, c.instance.InstanceID) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
		}
	}
}

func (c *ArthurHltClient) DeRegister() {
	c.ticker.Stop()
	c.client.UnregisterInstance(vars.AppName, c.instance.InstanceID)
}
