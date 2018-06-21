package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	sendSizeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sendSizeCounter",
			Help: "domainname send size counter",
		},
		[]string{"host"},
	)
	receiveSizeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "receiveSizeCounter",
			Help: "domainname receive size counter",
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
	sixCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "status6xxCounter",
			Help: "domainname 6xx counter",
		},
		[]string{"host"},
	)
)

func init() {
	prometheus.MustRegister(sendSizeCounter)
	prometheus.MustRegister(receiveSizeCounter)
	prometheus.MustRegister(zeroCounter)
	prometheus.MustRegister(twoCounter)
	prometheus.MustRegister(threeCounter)
	prometheus.MustRegister(fourCounter)
	prometheus.MustRegister(fiveCounter)
	prometheus.MustRegister(sixCounter)
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

type savefile struct {
	Host counters `json:"host"`
}

type counters struct {
	Zero            float64 `json:"zero"`
	Tow             float64 `json:"tow"`
	Three           float64 `json:"three"`
	Four            float64 `json:"four"`
	Five            float64 `json:"five"`
	Six             float64 `json:"six"`
	Send_counter    float64 `json:"send_counter"`
	Receive_counter float64 `json:"receive_counter"`
}

func getFile() (filepath string) {
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
	filepath = path + filename
	return
}

func tailLog(ch chan log) {
	filepath := getFile()
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

func addSaveDate(filepath string) {
	date, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("error of read save file")
	}
	var datemap map[string]counters
	json.Unmarshal(date, &datemap)
	for k, v := range datemap {
		sendSizeCounter.With(prometheus.Labels{"host": k}).Add(v.Send_counter)
		receiveSizeCounter.With(prometheus.Labels{"host": k}).Add(v.Receive_counter)
		zeroCounter.With(prometheus.Labels{"host": k}).Add(v.Zero)
		twoCounter.With(prometheus.Labels{"host": k}).Add(v.Tow)
		threeCounter.With(prometheus.Labels{"host": k}).Add(v.Three)
		fourCounter.With(prometheus.Labels{"host": k}).Add(v.Four)
		fiveCounter.With(prometheus.Labels{"host": k}).Add(v.Five)
		sixCounter.With(prometheus.Labels{"host": k}).Add(v.Six)
	}
}

func counter(ch chan log) {
	var path string = "/var/log/containers/"
	var file string = "ingress-log-counter.json"
	filepath := path + file
	var metrics map[string]*counters
	metrics = make(map[string]*counters)
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		addSaveDate(filepath)
	}
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			save, _ := json.Marshal(metrics)
			err := ioutil.WriteFile(filepath, save, 0664)
			if err != nil {
				fmt.Println("err or write save file", err.Error())
			}
		}
	}()
	for l := range ch {
		if metrics[l.Host] == nil {
			var p *counters = new(counters)
			metrics[l.Host] = p
		}
		sendsize, _ := strconv.ParseFloat(l.Sent_bytes_body, 64)
		receivesize, _ := strconv.ParseFloat(l.Request_length, 64)
		sendSizeCounter.With(prometheus.Labels{"host": l.Host}).Add(sendsize)
		metrics[l.Host].Send_counter = metrics[l.Host].Send_counter + sendsize
		receiveSizeCounter.With(prometheus.Labels{"host": l.Host}).Add(receivesize)
		metrics[l.Host].Receive_counter = metrics[l.Host].Receive_counter + receivesize
		s := string([]byte(l.Status)[:1])
		switch s {
		case "0":
			zeroCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
			metrics[l.Host].Zero++
		case "2":
			twoCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
			metrics[l.Host].Tow++
		case "3":
			threeCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
			metrics[l.Host].Three++
		case "4":
			fourCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
			metrics[l.Host].Four++
		case "5":
			fiveCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
			metrics[l.Host].Five++
		case "6":
			fiveCounter.With(prometheus.Labels{"host": l.Host}).Add(1)
			metrics[l.Host].Six++
		default:
			fmt.Println("Unknow http code", l.Status)
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
