package main

import (
	"fmt"
	//"net/http"
	//"regexp"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var promefile string = "/opt/metric/bond_speed.prom"

func getBondList() (ifname []string) {
	links, err := netlink.LinkList()
	if err != nil {
		fmt.Errorf("Get link list err:", err.Error())
	}
	for _, link := range links {
		if link.Type() == "bond" {
			ifname = append(ifname, link.Attrs().Name)
		}
	}
	return
}

func getBondSpeed(bond string) (speed string) {
	filename := path.Join("/sys/class/net/", bond, "speed")
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf("Error open bond speed file", err.Error())
	}
	speed = strings.Replace(string(buf), "\n", "", -1)
	return
}

func formatFileContent(bondList []string) (content string) {
	for _, bond := range bondList {
		var data string = `bond_speed{name="__name__"} __speed__`
		speed := getBondSpeed(bond)
		data = strings.Replace(data, "__name__", bond, 1)
		data = strings.Replace(data, "__speed__", speed, 1)
		content = content + data + "\n"
	}
	return
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, data)
// }

func main() {
	// http.HandleFunc("/", handler)
	// http.ListenAndServe("0.0.0.0:9001", nil)
	bondList := getBondList()
	content := formatFileContent(bondList)
	file, err := os.OpenFile(promefile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Errorf("Open/Create prometheus metric in /opt/metric file err :", err.Error())
	}
	defer file.Close()
	file.WriteString(content)
}
