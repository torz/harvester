package main

import (
	"fmt"
	"io/ioutil"
	//"os"
	"net/http"
	"net/url"
	"crypto/tls"
	"sync"
)

var wg sync.WaitGroup
/* 
   check stuff here http://httpbin.org/
   my ip http://httpbin.org/ip
   my user agent http://httpbin.org/user-agent
   proxies 202.195.244.211:8080 93.170.36.3:1080
   124.67.215.9	85.114.141.191
*/

type Harvester struct {
	proxy		string
	client		http.Client
	userAgent	string
}

func (h *Harvester) setClient() {

	urli :=url.URL{}
	urlProxy, _ := urli.Parse(h.proxy)

	transport := http.Transport{}
	transport.Proxy = http.ProxyURL(urlProxy)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	h.client = http.Client{}
	h.client.Transport = &transport
}

func (h *Harvester) get(myurl string) {

	myreq, _ := http.NewRequest("GET", myurl, nil)
	myreq.Header.Set("User-Agent", h.userAgent)

	resp, _ := h.client.Do(myreq)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", body)
}

func main() {
	mysites := []string{"http://httpbin.org/ip", "http://httpbin.org/user-agent", "http://httpbin.org/headers"}
	myproxy := []string{"http://183.224.1.30/", "http://221.130.23.150/", "http://221.130.23.144/"}
	myuas := []string{"chrome", "firefox", "ie"}

	for i := range myuas {
		wg.Add(1)
		go func(i int){
			h := Harvester{}
			h.proxy = myproxy[i]
			fmt.Println(h.proxy)
			h.userAgent = myuas[i]
			h.setClient()
			//fmt.Printf("%v\n", h.client.Transport)
			h.get(mysites[0])
			h.get(mysites[1])
			wg.Done()
		}(i);
	}
	wg.Wait()
}
