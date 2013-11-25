package gprkernel

import (
	"fmt"
	"net"
)

type Config struct {
	host string
	port float64
}

func Run(config *map[string] interface {}) {
	lConfig := Config{}
	if v, ok := (*config)["localAddr"].(string); ok {
		lConfig.host = v
	} else {
		lConfig.host = "0.0.0.0"
	}
	if v, ok := (*config)["localPort"].(float64); ok {
		lConfig.port = v
	}
	if v, ok := (*config)["routes"].(map[string]interface {}); ok {
		routes := map[string]Config{}
		for host, route := range v {
			if route, ok := route.(map[string]interface {}); ok {
				rConfig := Config{}
				if v, ok := route["dstAddr"].(string); ok {
					rConfig.host = v
				}
				if v, ok := route["dstPort"].(float64); ok {
					rConfig.port = v
				}
				routes[host] = rConfig
			}
		}
		go Router(&lConfig, &routes)
	} else {
		rConfig := Config{}
		if v, ok := (*config)["dstAddr"].(string); ok {
			rConfig.host = v
		}
		if v, ok := (*config)["dstPort"].(float64); ok {
			rConfig.port = v
		}
		go Proxy(&lConfig, &rConfig)
	}
}

func Proxy(lConfig *Config, rConfig *Config) {
	addressLocal := fmt.Sprint(*lConfig.host , ":" , *lConfig.port)
	tcpAddrLocal, err := net.ResolveTCPAddr("tcp4", addressLocal)
	if err != nil {fmt.Println(err)}

	listener, err := net.ListenTCP("tcp", tcpAddrLocal)
	if err != nil {fmt.Println(err)}

	addressDst := fmt.Sprint(*rConfig.host , ":" , *rConfig.port)
	tcpAddrDst, err := net.ResolveTCPAddr("tcp4", addressDst)
	if err != nil {fmt.Println(err)}

	for {
		connLocal, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func(connLocal *net.Conn) {
			connDst, err := net.DialTCP("tcp", nil, tcpAddrDst)
			if err != nil {fmt.Println(err)}
			fmt.Println(connDst)
		}(connLocal)
	}
}

func Router(lConfig *Config, routes *map[string]Config) {
}
