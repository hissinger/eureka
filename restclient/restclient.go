package restclient

import (
	"bytes"
	"encoding/json"
	"eureka/vars"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type RestClient struct {
	client     *http.Client
	instanceID string
	instance   map[string]interface{}
	ticker     *time.Ticker
	done       chan bool
}

func NewClient() *RestClient {
	port := strconv.Itoa(vars.LocalPort)
	baseURL := "http://" + vars.LocalIP + ":" + port

	c := &RestClient{instanceID: vars.LocalIP + ":" + vars.AppName + ":" + port}

	dir, _ := os.Getwd()
	data, _ := ioutil.ReadFile(dir + "/restclient/template.json")
	c.instance = map[string]interface{}{}
	json.Unmarshal(data, &c.instance)
	instance := c.instance["instance"].(map[string]interface{})
	instance["instanceId"] = c.instanceID
	instance["app"] = vars.AppName
	instance["hostName"] = vars.Hostname
	instance["ipAddr"] = vars.LocalIP
	instance["vipAddress"] = vars.AppName
	instance["port"] = port
	instance["status"] = "UP"
	instance["homePageUrl"] = baseURL + "/"
	instance["statusPageUrl"] = baseURL + "/info"
	instance["healthCheckUrl"] = baseURL + "/health"

	return c
}

func (c *RestClient) Register() {
	c.client = &http.Client{}

	data, _ := json.Marshal(c.instance)
	fmt.Println(string(data))
	req, _ := http.NewRequest(http.MethodPost, vars.EurekaServerURL+"/apps/"+vars.AppName, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", "localhost:8761")
	req.Header.Set("Connection", "Keep-Alive")

	err := c.sendHttpRequest(req)
	if err != nil {
		return
	}

	c.sendHeartBeat()
}

func (c *RestClient) sendHeartBeat() {
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
				c.heartbeat()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *RestClient) heartbeat() {
	req, _ := http.NewRequest(http.MethodPut, c.getURL(), nil)
	req.Header.Set("Host", "localhost:8761")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Transfer-Encoding", "chunked")

	err := c.sendHttpRequest(req)
	if err != nil {
		return
	}
}

func (c *RestClient) DeRegister() {
	c.ticker.Stop()

	fmt.Println("Trying to deregister application...")

	// update status
	instance := c.instance["instance"].(map[string]interface{})
	instance["status"] = "DOWN"
	data, _ := json.Marshal(c.instance)
	req, _ := http.NewRequest(http.MethodPost, vars.EurekaServerURL+"/apps/"+vars.AppName, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", "localhost:8761")
	req.Header.Set("Connection", "Keep-Alive")

	err := c.sendHttpRequest(req)
	if err != nil {
		return
	}

	// Deregister
	req, _ = http.NewRequest(http.MethodDelete, c.getURL(), nil)
	req.Header.Set("Host", "localhost:8761")
	req.Header.Set("Connection", "Keep-Alive")

	err = c.sendHttpRequest(req)
	if err != nil {
		return
	}

	fmt.Println("Deregistered application, exiting. Check Eureka...")
}

func (c *RestClient) sendHttpRequest(req *http.Request) error {
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(bodyBytes))

	return nil
}

func (c *RestClient) getURL() string {
	return vars.EurekaServerURL + "/apps/" + vars.AppName + "/" + c.instanceID
}
