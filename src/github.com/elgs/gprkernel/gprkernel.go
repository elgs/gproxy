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
	addressLocal := fmt.Sprint(lConfig.host , ":" , lConfig.port)
	tcpAddrLocal, err := net.ResolveTCPAddr("tcp4", addressLocal)
	if err != nil {fmt.Println(err)}

	listener, err := net.ListenTCP("tcp", tcpAddrLocal)
	if err != nil {fmt.Println(err)}

	addressDst := fmt.Sprint(rConfig.host , ":" , rConfig.port)
	tcpAddrDst, err := net.ResolveTCPAddr("tcp4", addressDst)
	if err != nil {fmt.Println(err)}

	for {
		connLocal, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		connDst, err := net.DialTCP("tcp", nil, tcpAddrDst)
		if err != nil {fmt.Println(err)}
		localChan := make(chan []byte, 10)
		dstChan := make(chan []byte, 10)
		go func() {
			for {
				select {
				case localChan <- func() ([]byte) {
					var buffer = make([]byte, 1024)
					n, err := connLocal.Read(buffer[0:])
					if err != nil {fmt.Println(err)}
					fmt.Println(1, n, string(buffer[0:n]))
					return buffer[0:n]
				}():
				case dstChan <- func() ([]byte) {
					var buffer = make([]byte, 1024)
					n, err := connDst.Read(buffer[0:])
					if err != nil {fmt.Println(err)}
					fmt.Println(2, n, string(buffer[0:n]))
					return buffer[0:n]
				}():
				}
			}
		}()

		go func() {
			for {
				select {
				case buffer := <-localChan:
					fmt.Println(3, string(buffer))
					connDst.Write(buffer)
				case buffer := <-dstChan:
					fmt.Println(4, string(buffer))
					connLocal.Write(buffer)
				}
			}
		}()
	}

}

func Router(lConfig *Config, routes *map[string]Config) {
}
