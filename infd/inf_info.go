package main

import (
	"fmt"
	"os/exec"
	//	"reflect"
	"strings"
)

func getIfname(ifnumber string) string {
	response, err := exec.Command("get_ifname.sh", ifnumber).Output()
	if err != nil {
		fmt.Println("fuck")
	}
	return strings.Replace(string(response), "\"", "", -1)
}

func main() {
	response, err := exec.Command("get_ip.sh").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	ip := string(response)
	for _, j := range strings.Split(ip, "\n") {
		for m, n := range strings.Split(j, " ") {
			if m%2 == 1 {
				fmt.Println(n)
				ifname := getIfname(n)
				fmt.Println(ifname)
			}
			//fmt.Println(n)
			//fmt.Println(reflect.TypeOf(m))
		}
		//fmt.Println(reflect.TypeOf(j))
	}
}
