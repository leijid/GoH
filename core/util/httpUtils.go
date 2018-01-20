package util

import (
	"net/http"
	"io/ioutil"
	"net"
	"time"
	"strings"
	"GoH/core/constant"
)

var client *http.Client

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*constant.HttpTimeOut)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * constant.HttpTimeOut))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * constant.HttpTimeOut,
		},
	}
}

func HttpGet(url string) string {
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func HttpPost(url string, header map[string]string, postData string) string {
	req, err := http.NewRequest("POST", "http://www.baidu.com", strings.NewReader(postData))
	if err != nil {
		panic(err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
