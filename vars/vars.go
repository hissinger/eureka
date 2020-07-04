package vars

import (
	"fmt"
	"net"
	"os"
)

var (
	EurekaServerURL   = "http://192.168.123.1:8761/eureka"
	LocalIP           = "127.0.0.1"
	Hostname          = "localhost"
	LocalPort         = 8080
	AppName           = "GO-EUREKA"
	HeartBeatInterval = 30
)

func init() {
	// get local ip
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	LocalIP = conn.LocalAddr().(*net.UDPAddr).IP.String()

	// get hostname
	Hostname, _ = os.Hostname()

	// addr, err := net.LookupIP(Hostname)
	// if err == nil {
	// 	for _, ip := range addr {
	// 		if ip.To4() != nil && ip.IsLoopback() == false {
	// 			LocalIP = addr[0].String()
	// 			break
	// 		}
	// 	}
	// }
	// fmt.Println(Hostname)
	// fmt.Println(LocalIP)
}
