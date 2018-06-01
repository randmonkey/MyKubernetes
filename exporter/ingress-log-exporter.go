package main

import (
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	sendSizeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sendSizeCounter",
			Help: "domainname send size counter",
		},
		[]string{"host"},
	)
	reciveSizeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reciveSizeCounter",
			Help: "domainname recive size counter",
		},
		[]string{"host"},
	)
	zeroCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status0xxCounter",
			Help: "domainname 0xx counter",
		},
		[]string{"host"},
	)
	twoCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status2xxCounter",
			Help: "domainname 2xx counter",
		},
		[]string{"host"},
	)
	threeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status3xxCounter",
			Help: "domainname 3xx counter",
		},
		[]string{"host"},
	)
	fourCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status4xxCounter",
			Help: "domainname 4xx counter",
		},
		[]string{"host"},
	)
	fiveCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status5xxCounter",
			Help: "domainname 5xx counter",
		},
		[]string{"host"},
	)
)

func init() {
	prometheus.MustRegister(sendSizeCounter)
	prometheus.MustRegister(reciveSizeCounter)
	prometheus.MustRegister(zeroCounter)
	prometheus.MustRegister(twoCounter)
	prometheus.MustRegister(threeCounter)
	prometheus.MustRegister(fourCounter)
	prometheus.MustRegister(fiveCounter)
}

type logstr struct {
	Log string `json:"log"`
}

type log struct {
	Status          string `json:"status"`
	Sent_bytes      string `json:"bytes_sent"`
	Sent_bytes_body string `json:"body_bytes_sent"`
	Request_length  string `json:"request_length"`
	Host            string `json:"host"`
}

type counters struct {
	tow            float64
	three          float64
	four           float64
	five           float64
	send_counter   float64
	recive_counter float64
}

func tailLog(ch chan log) {
	var filename string
	var path string = "/var/log/containers/"
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("error of read log directory")
	}
	for _, v := range dir {
		if j, _ := regexp.MatchString("nginx-ingress-controller-(.*)log", v.Name()); j {
			filename = v.Name()
		}
	}
	filepath := path + filename
	t, err := tail.TailFile(filepath, tail.Config{Follow: true})
	if err != nil {
		fmt.Println("error of tail access.log!", err.Error())
	}
	for line := range t.Lines {
		var ls logstr
		var logs log
		err := json.Unmarshal([]byte(line.Text), &ls)
		if err != nil {
			fmt.Println("error of unmarshal line.text -> logstr", err.Error())
		}
		ls.Log = strings.Replace(ls.Log, "\\", "", -1)
		if isok, _ := regexp.MatchString("^[WIE]?([0-9]+)", ls.Log); !isok {
			err := json.Unmarshal([]byte(ls.Log), &logs)
			if err != nil {
				fmt.Println("error of unmarshal logstr -> log", err.Error())
			}
			ch <- logs
		}
	}
}

func counter(ch chan log) {
	for l := range ch {
		sendsize, _ := strconv.ParseFloat(l.Sent_bytes, 64)
		recivesize, _ := strconv.ParseFloat(l.Request_length, 64)
		sendSizeCounter.With(prometheus.Labels{"host": l.Host}).Add(sendsize)
		reciveSizeCounter.With(prometheus.Labels{"host": l.Host}).Add(recivesize)
		s := string([]byte(l.Status)[:1])
		switch s {
		case "0":
			zeroCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
		case "2":
			twoCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
		case "3":
			threeCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
		case "4":
			fourCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
		case "5":
			fiveCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
		default:
			fmt.Println("error http code", l.Status)
		}
	}
}

func main() {
	ch := make(chan log, 1000)
	go tailLog(ch)
	go counter(ch)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8888", nil)
}
