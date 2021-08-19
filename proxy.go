package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	proxies    []string
	lineCount  int    = 0
	active     int    = 0
	colorGreen string = "\033[32m"
	colorReset string = "\033[0m"
	colorRed   string = "\033[31m"
)

func main() {
	fmt.Println("[" + colorGreen + "+" + colorReset + "] " + "Proxy checker by rossSec ")
	readFile()
	for _, proxy := range proxies {
		go checkProxy(proxy)
	}
	fmt.Println("Press CTRL+C to Exit" + "\n")
	time.Sleep(5000 * time.Second)
}

func readFile() {
	file, err := os.Open("proxies.txt")
	if err != nil {
		fmt.Println("[" + colorRed + "-" + colorReset + "] " + "Error reading proxies.txt")
		os.Exit(0)
	}
	fmt.Println("[" + colorGreen + "+" + colorReset + "] " + "Reading proxies.txt" + "\n")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
		proxies = append(proxies, scanner.Text())
	}
}

func writeFile(proxy string) {
	f, err := os.Create("checked.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(proxy)
	f.Close()
}

func checkProxy(proxy string) {
	proxyStr := "http://" + proxy
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}
	urlStr := "https://google.com/"
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
	}
	response, err := client.Do(request)
	if err != nil {
		return
	} else {
		active++
		counter := fmt.Sprintf("["+colorGreen+"Checking"+colorReset+"] "+"%d"+"/"+"%d", active, lineCount)
		fmt.Print(counter, "\r")
		writeFile(proxy)
	}
	defer response.Body.Close()
}
