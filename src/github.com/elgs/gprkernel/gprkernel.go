package gprkernel

import "fmt"

type Config struct {
	host string
	port float64
	routes map[string]interface {}
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
		Router(&lConfig, &routes)
	} else {
		rConfig := Config{}
		if v, ok := (*config)["dstAddr"].(string); ok {
			rConfig.host = v
		}
		if v, ok := (*config)["dstPort"].(float64); ok {
			rConfig.port = v
		}
		Proxy(&lConfig, &rConfig)
	}
}

func Proxy(lConfig *Config, rConfig *Config) {
	fmt.Println("In Proxy")
}

func Router(lConfig *Config, routes *map[string]Config) {
	fmt.Println("In Router")
}
