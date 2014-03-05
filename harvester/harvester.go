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
}

func (h *Harvester) setClient() {

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

func (h *Harvester) get(myurl string) {

	myreq, _ := http.NewRequest("GET", myurl, nil)
	myreq.Header.Set("User-Agent", h.userAgent)

	resp, err := h.client.Do(myreq)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("%s\n", body)
}

func main() {
	mysites := []string{"http://httpbin.org/ip", "http://httpbin.org/user-agent", "http://httpbin.org/headers"}
	myproxy := []string{"http://183.224.1.30/", "http://221.130.23.150/", "http://221.130.23.144/", "http://94.205.181.212/" ,"http://173.201.95.24/"}
	myuas := []string{"chrome", "firefox", "ie", "opera", "safari"}

	for i := range myproxy {
		wg.Add(1)
		go func(i int){
			h := Harvester{}
			h.proxy = myproxy[i]
			h.userAgent = myuas[i]
			h.setClient()
			h.get(mysites[0])
			h.get(mysites[1])
			wg.Done()
		}(i)
	}
	wg.Wait()
}
