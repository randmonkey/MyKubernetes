package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var promefile string = "/opt/metric/bond_speed.prom"

func getBondList() (ifname []string, err error) {
	links, err := netlink.LinkList()
	if err != nil {
		err = fmt.Errorf("get link error ", err.Error())
		return
	}
	for _, link := range links {
		if link.Type() == "bond" {
			ifname = append(ifname, link.Attrs().Name)
		}
	}
	return
}

func getBondSpeed(bond string) (speed string, err error) {
	filename := path.Join("/sys/class/net/", bond, "speed")
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("Error open bond speed file", err.Error())
		return
	}
	speed = strings.Replace(string(buf), "\n", "", -1)
	return
}

func formatFileContent(bondList []string) (content string, err error) {
	for _, bond := range bondList {
		var data string = `bond_speed{name="__name__"} __speed__`
		speed, e := getBondSpeed(bond)
		if e != nil {
			err = fmt.Errorf("get bond speed err", e.Error())
			return
		}
		data = strings.Replace(data, "__name__", bond, 1)
		data = strings.Replace(data, "__speed__", speed, 1)
		content = content + data + "\n"
	}
	return
}

func main() {
	_, err := os.Stat("/opt/metric")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("/opt/metric", 0777)
		}
	}
	bondList, _ := getBondList()
	content, _ := formatFileContent(bondList)
	file, err := os.OpenFile(promefile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("can't open file prome", err)
	}
	defer file.Close()
	file.WriteString(content)
}
