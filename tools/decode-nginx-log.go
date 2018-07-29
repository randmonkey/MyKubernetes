package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"bufio"
	"os"
)

type logstr struct {
	Log string `json:"log"`
}

func main() {
	var ls logstr
	file, err := os.Open("fail.json")
	if err != nil {
		fmt.Println("error of open log file", err.Error())
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		line := reader.Text()
		_ = json.Unmarshal([]byte(line), &ls)
		fmt.Println(ls.Log)
	}
}
