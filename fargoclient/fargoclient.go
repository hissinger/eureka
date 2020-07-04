package fargoclient

import (
	"eureka/vars"
	"fmt"
	"strconv"
	"time"

	"github.com/hudl/fargo"
)

type FargoClient struct {
	instanceID string
	conn       fargo.EurekaConnection
	instance   *fargo.Instance
	ticker     *time.Ticker
	done       chan bool
}

func NewClient() *FargoClient {
	return &FargoClient{instanceID: vars.LocalIP + ":" + vars.AppName + ":" + strconv.Itoa(vars.LocalPort)}
}

func (c *FargoClient) Register() {
	conf := fargo.Config{}
	conf.Eureka.ServiceUrls = []string{vars.EurekaServerURL}
	c.conn = fargo.NewConnFromConfig(conf)
	c.conn.UseJson = true

	c.instance = &fargo.Instance{
		InstanceId:        c.instanceID,
		HostName:          vars.Hostname,
		Port:              vars.LocalPort,
		PortEnabled:       true,
		SecurePort:        443,
		SecurePortEnabled: true,
		App:               vars.AppName,
		IPAddr:            vars.LocalIP,
		VipAddress:        vars.AppName,
		SecureVipAddress:  vars.AppName,
		CountryId:         1,
		Status:            fargo.UP,
		DataCenterInfo:    fargo.DataCenterInfo{Name: fargo.MyOwn},
		HomePageUrl:       "http://" + vars.LocalIP + ":" + strconv.Itoa(vars.LocalPort) + "/",
		StatusPageUrl:     "http://" + vars.LocalIP + ":" + strconv.Itoa(vars.LocalPort) + "/actuator/info",
		HealthCheckUrl:    "http://" + vars.LocalIP + ":" + strconv.Itoa(vars.LocalPort) + "/actuator/health",
		LeaseInfo: fargo.LeaseInfo{
			RenewalIntervalInSecs: 30,
			DurationInSecs:        90,
		},
	}

	c.instance.SetMetadataString("management.port", strconv.Itoa(vars.LocalPort))

	c.conn.RegisterInstance(c.instance)

	c.sendHeartBeat()
}

func (c *FargoClient) sendHeartBeat() {
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
				c.conn.HeartBeatInstance(c.instance)
			case <-c.done:
				return
			}
		}
	}()
}

func (c *FargoClient) DeRegister() {
	close(c.done)

	c.instance.Status = fargo.DOWN
	c.conn.RegisterInstance(c.instance)
	c.conn.DeregisterInstance(c.instance)
}
