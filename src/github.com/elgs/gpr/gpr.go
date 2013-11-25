package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/elgs/gprkernel"
	"os"
)

func main() {
	input := args()
	for i := range input {
		start(input[i])
	}
	select {}
}

func start(configFile string) {
	configs, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var f interface{}
	err = json.Unmarshal(configs, &f)
	if err != nil {fmt.Println(err)}
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
