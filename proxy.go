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
	fmt.Println("[" + colorGreen + "+" + colorReset + "] " + "Proxy Checker," + colorGreen + " https://github.com/rossSec/asyncChecker" + colorReset)
	readFile()
	fmt.Println("Press CTRL+C to Exit" + "\n")
	time.Sleep(3 * time.Second)
	for _, proxy := range proxies {
		go checkProxy(proxy)
	}
	time.Sleep(500000 * time.Second)
}

func readFile() {
	file, err := os.Open("proxies.txt")
	if err != nil {
		fmt.Println("[" + colorRed + "-" + colorReset + "] " + "Error reading proxies.txt")
		os.Exit(0)
	}
	fmt.Println("[" + colorGreen + "+" + colorReset + "] " + "Reading proxies.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
		proxies = append(proxies, scanner.Text())
	}
	counter := fmt.Sprintf("["+colorGreen+"+"+colorReset+"] "+"Loaded "+colorGreen+"%d"+colorReset+" proxies", lineCount)
	fmt.Println(counter + "\n")

}

func writeFile(proxy string) {
	f, err := os.OpenFile("checked.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(proxy + "\n"); err != nil {
		panic(err)
	}
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
		fmt.Println("[" + colorRed + "-" + colorReset + "] " + colorRed + proxy + colorReset)
		return
	} else {
		fmt.Println("[" + colorGreen + "+" + colorReset + "] " + colorGreen + proxy + colorReset)
		writeFile(proxy)
	}
	defer response.Body.Close()
}
