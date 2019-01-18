package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	certutil "k8s.io/client-go/util/cert"
)

var ErrNotInCluster = errors.New("unable to load in-cluster config")

type Config struct {
	Host     string
	APIPath  string
	Username string
	Password string
	TLSClientConfig
	BearerToken     string
	BearerTokenFile string
}

type TLSClientConfig struct {
	Insecure   bool
	ServerName string
	CertFile   string
	KeyFile    string
	CAFile     string
	CertData   []byte
	KeyData    []byte
	CAData     []byte
}

func getInClusterConfig() (*Config, error) {
	const (
		tokenFile  = "var/run/secrets/kubernetes.io/serviceaccount/token"
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, ErrNotInCluster
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	tlsClientConfig := TLSClientConfig{}

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		fmt.Printf("Expected to load root CA config from %s, but got err : %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &Config{
		Host:            "http://" + host + ":" + port,
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}, nil
}

func main() {
	fmt.Println("get k8s config")
	c, err := getInClusterConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}
