package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"crypto/tls"
	"sync"
	"time"
)

/* 
   check stuff here http://httpbin.org/
   my ip http://httpbin.org/ip
   my user agent http://httpbin.org/user-agent
*/

var wg sync.WaitGroup

type Harvester struct {
	proxy		string
	client		http.Client
	userAgent	string
	//resp		[]byte
}

func (h *Harvester) setClientProxy() {

	urli :=url.URL{}
	urlProxy, _ := urli.Parse(h.proxy)

	timeout := time.Duration(8 * time.Second)

	transport := http.Transport{}
	transport.Proxy = http.ProxyURL(urlProxy)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	transport.ResponseHeaderTimeout = timeout

	h.client = http.Client{}
	h.client.Transport = &transport
}

func (h *Harvester) get(myurl string) (body []byte) {

	myreq, _ := http.NewRequest("GET", myurl, nil)
	myreq.Header.Set("User-Agent", h.userAgent)

	resp, err := h.client.Do(myreq)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	//fmt.Printf("%s\n", body)
	return
}

func main() {
	datasource := "http://localhost:8888/urls"

	mysites := []string{"http://httpbin.org/ip", "http://httpbin.org/user-agent", "http://httpbin.org/headers"}
	myproxy := []string{"http://183.224.1.30/", "http://221.130.23.150/", "http://221.130.23.144/", "http://94.205.181.212/"}
	myuas := []string{"chrome", "firefox", "ie", "opera", "safari"}

	jsonharvester := Harvester{}
	fmt.Printf("%s\n", jsonharvester.get(datasource))

	for i := range myproxy {
		wg.Add(1)
		go func(i int){
			h := Harvester{}
			h.proxy = myproxy[i]
			h.userAgent = myuas[i]
			h.setClientProxy()
			fmt.Printf("%s\n", h.get(mysites[0]))
			fmt.Printf("%s\n", h.get(mysites[1]))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
