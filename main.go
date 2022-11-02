package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	// File For urls
	Filename string `short:"f" long:"file" description:"subdomains file" required:"true"`
	// Concurrency For the requests
	Concurrency int `short:"c" long:"concurrency" default:"25" description:"Concurrency For Requests"`
}

var trans = &http.Transport{
	MaxIdleConns:      30,
	IdleConnTimeout:   time.Second,
	DisableKeepAlives: true,
	// Skip Certificate Error
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	TLSHandshakeTimeout: 5 * time.Second,
	Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: time.Second,
	}).Dial,
}

var client = &http.Client{
	// Passing transport var
	Transport: trans,
	// prevent follow redirect
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Timeout: 5 * time.Second,
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println("try to use -h to show the options and usage :)")
		// Exit from script after print this
		return
	}
	if !isExists(opts.Filename) {
		fmt.Println("File Not Found")
		return
	}
	data := readData(opts.Filename)
	var wg sync.WaitGroup
	for i := 0; i < opts.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range data {
				match := randomMatcher()
				putFile(d, match)
				checkIfVuln(d, match)
			}
		}()
	}

	wg.Wait()

}

func readData(filename string) <-chan string {

	urls := make(chan string)

	file, _ := os.Open(filename)

	data := bufio.NewScanner(file)

	go func() {
		defer file.Close()
		defer close(urls)
		for data.Scan() {
			url := strings.ToLower(data.Text())
			if !strings.Contains(url, "https://") && !strings.Contains(url, "http://") {
				url = "https://" + url
			}
			urls <- url
		}

	}()
	return urls
}

func isExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func putFile(url string, match string) {
	req, err := http.NewRequest(http.MethodOptions, url, nil)

	if err != nil {
		fmt.Println(err, "Error")
		return
	}
	req.Header.Set("Connection", "close")
	// Do the request
	resp, err := client.Do(req)

	if err != nil {
		// Check if host is up or not
		fmt.Println("error in host")
		return
	}
	defer resp.Body.Close()
	method := strings.Contains(strings.ToLower(resp.Header.Get("Allow")), "put")

	if !method {
		req.Header.Set("X-HTTP-Method-Override", "PUT")
	}

	req2, err := http.NewRequest(http.MethodPut, url+"/snowman.txt", strings.NewReader("Put method enabled,Uploaded by Mrsnowman! "+match))

	if err != nil {
		fmt.Println("Opps there is error")
		return
	}

	resp2, err := client.Do(req2)

	if err != nil {
		fmt.Println("Opps there is error")
		return
	}

	defer resp2.Body.Close()
}

func checkIfVuln(url string, match string) (vuln bool) {

	vuln = false
	req2, err := http.NewRequest(http.MethodGet, url+"/snowman.txt", nil)

	if err != nil {
		fmt.Println("Opps there is error", err)
		return
	}

	resp2, err := client.Do(req2)

	if err != nil {
		fmt.Println("Opps there is error", err)
		return
	}

	body, _ := ioutil.ReadAll(resp2.Body)
	if strings.Contains(string(body), match) {
		vuln = true
		fmt.Println(url + " ->>> Vulnerable")
	}

	defer resp2.Body.Close()

	return vuln
}

func randomMatcher() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 40)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
