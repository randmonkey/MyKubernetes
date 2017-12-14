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

func getIfnameIp() map[string]string {
	var ifname_ip map[string]string
	ifname_ip = make(map[string]string)
	response, err := exec.Command("get_ip.sh").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	ip := string(response)
	for _, j := range strings.Split(ip, "\n") {
		for m, n := range strings.Split(j, " ") {
			if m%2 == 1 {
				ifnames := getIfname(n)
				ifname := strings.Replace(ifnames, "\n", "", -1)
				ifname_ip[ifname] = strings.Split(j, " ")[0]
			}
		}
	}
	return ifname_ip
}

func main() {
	ifip := getIfnameIp()
	fmt.Println(ifip)
}
