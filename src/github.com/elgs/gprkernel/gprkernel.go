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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	addressLocal := fmt.Sprint(lConfig.host , ":" , lConfig.port)
	tcpAddrLocal, err := net.ResolveTCPAddr("tcp4", addressLocal)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddrLocal)
	if err != nil {
		panic(err)
	}

	addressDst := fmt.Sprint(rConfig.host , ":" , rConfig.port)
	tcpAddrDst, err := net.ResolveTCPAddr("tcp4", addressDst)
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
			connDst, err := net.DialTCP("tcp", nil, tcpAddrDst)
			defer func() {
				connDst.Close()
				fmt.Println("Server connection closed.")
				connLocal.Close()
				fmt.Println("Client connection closed.")
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			if err != nil {
				panic(err)
			}
			go func() {
				defer func() {
					connDst.Close()
					fmt.Println("Server connection closed.")
					connLocal.Close()
					fmt.Println("Client connection closed.")
					if err := recover(); err != nil {
						fmt.Println(err)
					}
				}()
				var buffer = make([]byte, 4096)
				for {
					runtime.Gosched()
					n, err := connLocal.Read(buffer)
					if err != nil {
						panic(err)
					}
					if n > 0 {
						_, err := connDst.Write(buffer[0:n])
						if err != nil {
							panic(err)
						}
					}
				}
			}()

			var buffer = make([]byte, 4096)
			for {
				runtime.Gosched()
				n, err := connDst.Read(buffer)
				if err != nil {
					panic(err)
				}
				if n > 0 {
					_, err := connLocal.Write(buffer[0:n])
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
}

func Router(lConfig *Config, routes *map[string]Config) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	addressLocal := fmt.Sprint(lConfig.host , ":" , lConfig.port)
	tcpAddrLocal, err := net.ResolveTCPAddr("tcp4", addressLocal)
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
				connLocal.Close()
				fmt.Println("Client connection closed.")
				if err := recover(); err != nil {
					fmt.Println(err)
				}
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
					addressDst := fmt.Sprint(rConfig.host , ":" , rConfig.port)
					tcpAddrDst, err := net.ResolveTCPAddr("tcp4", addressDst)
					if err != nil {
						panic(err)
					}
					connDst, err := net.DialTCP("tcp", nil, tcpAddrDst)
					defer func() {
						connDst.Close()
						fmt.Println("Server connection closed.")
						if err := recover(); err != nil {
							fmt.Println(err)
						}
					}()
					if err != nil {
						panic(err)
					}
					_, err = connDst.Write(peep[0:n])
					if err != nil {
						panic(err)
					}
					go func() {
						defer func() {
							connDst.Close()
							fmt.Println("Server connection closed.")
							connLocal.Close()
							fmt.Println("Client connection closed.")
							if err := recover(); err != nil {
								fmt.Println(err)
							}
						}()
						var buffer = make([]byte, 4096)
						for {
							runtime.Gosched()
							n, err := connLocal.Read(buffer)
							if err != nil {
								panic(err)
							}
							if n > 0 {
								_, err := connDst.Write(buffer[0:n])
								if err != nil {
									panic(err)
								}
							}
						}
					}()

					var buffer = make([]byte, 4096)
					for {
						runtime.Gosched()
						n, err := connDst.Read(buffer)
						if err != nil {
							panic(err)
						}
						if n > 0 {
							_, err := connLocal.Write(buffer[0:n])
							if err != nil {
								panic(err)
							}
						}
					}
					break;
				}
			}
		}()
	}
}
