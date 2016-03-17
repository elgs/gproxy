package gprkernel

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type Config struct {
	host string
	port float64
}

func Run(config *map[string]interface{}) {
	lConfig := Config{}
	if v, ok := (*config)["localAddr"].(string); ok {
		lConfig.host = v
	} else {
		lConfig.host = "[::]"
	}
	if v, ok := (*config)["localPort"].(float64); ok {
		lConfig.port = v
	}
	if v, ok := (*config)["routes"].(map[string]interface{}); ok {
		routes := map[string]Config{}
		for host, route := range v {
			if route, ok := route.(map[string]interface{}); ok {
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

func pipe(connLocal *net.Conn, connDst *net.Conn, bufSize int) {
	var buffer = make([]byte, bufSize)
	for {
		runtime.Gosched()
		n, err := (*connLocal).Read(buffer)
		if err != nil {
			(*connLocal).Close()
			(*connDst).Close()
			break
		}
		if n > 0 {
			_, err := (*connDst).Write(buffer[0:n])
			if err != nil {
				(*connLocal).Close()
				(*connDst).Close()
				break
			}
		}
	}
}

func Proxy(lConfig *Config, rConfig *Config) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	addressLocal := fmt.Sprint(lConfig.host, ":", lConfig.port)
	tcpAddrLocal, err := net.ResolveTCPAddr("tcp", addressLocal)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddrLocal)
	if err != nil {
		panic(err)
	}

	addressDst := fmt.Sprint(rConfig.host, ":", rConfig.port)

	for {
		connLocal, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			connDst, err := net.Dial("tcp", addressDst)
			if err != nil {
				fmt.Println(err)
				connLocal.Close()
				fmt.Println("Client connection closed.")
				return
			}
			go pipe(&connLocal, &connDst, 4096)
			pipe(&connDst, &connLocal, 4096)
		}()
	}
}

func Router(lConfig *Config, routes *map[string]Config) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	addressLocal := fmt.Sprint(lConfig.host, ":", lConfig.port)
	tcpAddrLocal, err := net.ResolveTCPAddr("tcp", addressLocal)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddrLocal)
	if err != nil {
		panic(err)
	}
	for {
		connLocal, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
				connLocal.Close()
				fmt.Println("Client connection closed.")
			}()
			var peep = make([]byte, 4096)
			n, err := connLocal.Read(peep)
			if err != nil {
				panic(err)
			}
			headers := strings.Split(string(peep[0:n]), "\n")
			for _, header := range headers {
				if strings.HasPrefix(strings.ToLower(header), "host") {
					hostData := strings.Split(header, ":")
					route := strings.TrimSpace(hostData[1])
					rConfig, ok := (*routes)[route]
					if !ok {
						rConfig, ok = (*routes)["default"]
						if !ok {
							panic("No route found.")
						}
					}
					addressDst := fmt.Sprint(rConfig.host, ":", rConfig.port)
					connDst, err := net.Dial("tcp", addressDst)
					if err != nil {
						fmt.Println(err)
						connLocal.Close()
						fmt.Println("Client connection closed.")
						return
					}
					_, err = connDst.Write(peep[0:n])
					if err != nil {
						panic(err)
					}
					go pipe(&connLocal, &connDst, 4096)
					pipe(&connDst, &connLocal, 4096)
					break
				}
			}
		}()
	}
}
