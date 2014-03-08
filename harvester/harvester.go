package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"crypto/tls"
	"sync"
	"time"
	"math/rand"
)

/* 
   check stuff here http://httpbin.org/
   my ip http://httpbin.org/ip
   my user agent http://httpbin.org/user-agent
*/

var wg sync.WaitGroup

type Harvester struct {
	client			http.Client
	userAgent		string
	proxyAddress	string
	url				string
}

func (h *Harvester) setClientProxy(proxyAddress string) {

	urli :=url.URL{}
	h.proxyAddress = proxyAddress
	urlProxy, _ := urli.Parse(h.proxyAddress)

	timeout := time.Duration(5 * time.Second)

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
	return
}

func main() {
	datasource := "http://localhost:5000/urls"

	mysites := []string{"http://httpbin.org/ip", "http://httpbin.org/user-agent", "http://httpbin.org/headers"}
	myproxy := []string{"http://183.224.1.30/", "http://221.130.23.150/", "http://221.130.23.144/"}
	myuas := []string{"chrome", "firefox", "ie", "opera", "safari"}

	jsonharvester := Harvester{}
	fmt.Printf("%s\n", jsonharvester.get(datasource))

	for i := range mysites {
		wg.Add(1)
		go func(i int){
			h := Harvester{}
			x := rand.Intn(len(myuas))
			fmt.Println("x is: ", x)
			h.userAgent = myuas[x]
			y := rand.Intn(len(myproxy))
			fmt.Println("pr length: ", len(myproxy))
			fmt.Println("y is: ", y)
			h.setClientProxy(myproxy[y])
			fmt.Printf("%s\n", h.get(mysites[0]))
			fmt.Printf("%s\n", h.get(mysites[1]))
			fmt.Printf("%s\n", h.get(mysites[2]))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
