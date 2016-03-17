gpr
======

A go proxy server.
------
Supports TCP in proxy mode, HTTP and Websocket in router mode.

Supported Platforms
------
* Windows 32/64
* Linux 32/64/ARM
* OS X

Installation
------
`go get github.com/elgs/gproxy/gpr`

Startup (as command line utility)
------
* Start the proxy server with default configuration `gpr.json`:
`gpr`

* Start the proxy server with selected configuration `google.json`:
`gpr google.json`

* Start the proxy server with multiple configurations `google.json ms.json`:
`gpr google.json ms.json`

Use as go module
------
```go
package main

import (
	"encoding/json"
	"github.com/elgs/gprkernel"
	"time"
)

func main() {
	configProxy1 := `{
		"dstPort": 80,
		"localPort": 3000,
		"dstAddr": "www.microsoft.com"
	}`

	configProxy2 := `{
		"dstPort": 80,
		"localPort": 2000,
		"dstAddr": "www.google.com"
	}`

	configRouter := `{
		"localPort": 8000,
		"routes": {
			"hosta": {
				"dstAddr": "[::]",
				"dstPort": 10309
			},
			"default": {
				"dstAddr": "127.0.0.1",
				"dstPort": 10310
			}
		}
	}`

	var f1 interface{}
	json.Unmarshal([]byte(configProxy1), &f1)
	m1 := f1.(map[string]interface{})
	gprkernel.Run(&m1)

	var f2 interface{}
	json.Unmarshal([]byte(configProxy2), &f2)
	m2 := f2.(map[string]interface{})
	gprkernel.Run(&m2)

	var f3 interface{}
	json.Unmarshal([]byte(configRouter), &f3)
	m3 := f3.(map[string]interface{})
	gprkernel.Run(&m3)

	for {
		time.Sleep(time.Hour)
	}
}
```

Proxy mode
------
A configuration `gpr.json` looks like this:
```js
{
  "microsoft" : {
    "dstPort" : 80,
    "localPort" : 3000,
    "dstAddr" : "www.microsoft.com"
  },
  "google" : {
    "dstPort" : 80,
    "localPort" : 2000,
    "localAddr" : "[::]",
    "dstAddr" : "www.google.com"
  }
}
```
means that:
when the clients connect to 127.0.0.1:3000, they connect to `www.microsoft.com:80`, and when the clients connect to `127.0.0.1:2000`, they connect to `www.google.com:80`.

Please note that `localAddr` is not necessary, when omitted, the server will listen on all network interfaces.

Router mode
------
Router mode works only with HTTP, not even HTTPS. Proxy mode and router mode can be working together happily.
A configuration `gpr.json` looks like this:
```js
{
  "google_ms" : {
    "localPort" : 4000,
    "routes" : {
      "hostname_a" : {
        "dstAddr" : "www.microsoft.com",
        "dstPort" : 80
      },
      "hostname_b" : {
        "dstAddr" : "www.google.com",
        "dstPort" : 80
      },
      "ipv6" : {
        "dstPort" : 22,
        "localPort" : 4001,
        "dstAddr" : "[2607:1260:1234::1:0]"
      }
    }
  }
}
```
means that:
if multiple host names / domain names are bound to the proxy server, let's say `hostname_a` and `hostname_b`. When the clients connect to `hostname_a:4000`, they connect to `www.microsoft.com:80`, and when the clients connect to `hostname_b:4000`, they connect to `www.google.com:80`. If the clients connect to a host name which is not in the route table, `127.0.0.1:4000` from the proxy server itelf, for example, they connect to the default route `www.yahoo.com:80`.
