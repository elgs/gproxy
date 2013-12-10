package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/elgs/gprkernel"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("CPUs: ", runtime.NumCPU())
	input := args()
	for i := range input {
		start(input[i])
	}
	for {
		time.Sleep(time.Hour)
	}
}

func start(configFile string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	configs, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprint(configFile, " not found"))
	}

	var f interface{}
	err = json.Unmarshal(configs, &f)
	if err != nil {
		panic(err)
	}
	m := f.(map[string]interface{})
	for _, v := range m {
		if value, ok := v.(map[string]interface {}); ok {
			gprkernel.Run(&value)
		}
	}
}

func args() ([]string) {
	ret := []string{}
	if len(os.Args) <= 1 {
		ret = append(ret, "gpr.json")
	} else {
		for i := 1; i < len(os.Args); i++ {
			ret = append(ret, os.Args[i])
		}
	}
	return ret
}

func test() {
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
				"dstAddr": "127.0.0.1",
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
