package infd

import (
	"fmt"
	"os/exec"
	"strings"
)

func getIfname(ifnumber string) (ifname string, err error) {
	response, err := exec.Command("get_ifname.sh", ifnumber).Output()
	if err != nil {
		err = fmt.Errorf("Error invoke get_ifname.sh:", err)
		return
	}
	ifname := strings.Replace(string(response), "\"", "", -1)
	return
}

func getIfnameIp() (ifname_ip map[string]string, err error) {
	var ifname_ip map[string]string
	ifname_ip = make(map[string]string)
	response, err := exec.Command("get_ip.sh").Output()
	if err != nil {
		fmt.Errorf("Error invoke get_ip.sh:", err)
		return
	}
	ip := string(response)
	for _, j := range strings.Split(ip, "\n") {
		for m, n := range strings.Split(j, " ") {
			if m%2 == 1 {
				ifnames, err := getIfname(n)
				if err != nil {
					err = fmt.Errorf("Error invoke getIfname :", err)
					return
				}
				ifname := strings.Replace(ifnames, "\n", "", -1)
				ifname_ip[ifname] = strings.Split(j, " ")[0]
			}
		}
	}
	return
}
