package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func guide() {
	// do -h parser job
	fmt.Println("subChie is a tool for finding the subdomains of a web site")
	fmt.Println(`
				  ||
-h			shows this\/ 		- - -
-d			domain			-d example.com
-w			wordlist		-w Path
-g			google hakcing  -g

	`)
}

func banner() {
	fmt.Println(`
 ____        _       _ __ _     |_|   
/ ___| _   _| |__  /  _ _| | 	__  ____
\___ \| | | | '_ \| |    | |_ _|  /  ___\
 ___) | |_| | |_) | \ _ _|  __ |  | /___|
|____/ \__,_|_.__/ \ __ _|_| |_|  | \
 created by Bl00dBlu35
 - - - - - - - - - - - - - - - - - - - - - -`)
}

func create_file() (*os.File, error) {
	// creating a new file name output.txt
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return f, nil
}

func open_file(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return file, nil
}

func create_url(protocol string, scanner *bufio.Scanner, domain string) string {
	// string builder object
	var SB_url strings.Builder
	// url := "https://" + scanner.Text() + domain
	SB_url.WriteString(protocol)
	SB_url.WriteString("://")
	SB_url.WriteString(scanner.Text())
	SB_url.WriteString(".")
	SB_url.WriteString(domain)
	url := SB_url.String()
	return url
}

func check_subs(domain string, wordList string, bad_status_code [10]int) {
	// create the output.txt file
	output, _ := create_file()

	// opening the wordlist file
	wordListOb, _ := open_file(wordList)

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(wordListOb)

	//Read each line and print it to the console
	for scanner.Scan() {

		url := create_url("https", scanner, domain)
		fmt.Println(url)

		// the response of the url
		res, err := http.Head(url)
		if err != nil {
			continue
		}

		// search for status codes of res in bad status code and create a Output.txt file of urls with good status codes
		var flag bool = true
		for _, bad_status := range bad_status_code {
			if res.StatusCode == bad_status {
				flag = false
				break
			}
		}
		if flag {
			color.New(color.FgGreen, color.Bold).Println(strconv.Itoa(res.StatusCode) + " " + url)
			if _, err = output.WriteString(url + "\n"); err != nil {
				panic(err)
			}
		}

	}
}

func dorking(domain string) {
	query := fmt.Sprint("inurl:*.", domain)
	for page := 0; page <= 50; page++ {
		dork := fmt.Sprintf("https://www.google.com/search?q=%s+-www&filter=0&start=%d", query, page)

		fmt.Println(dork)

		output, _ := create_file()

		resp, err := http.Get(dork)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(body))
		// Parse response body for URLs using regular expressions
		urlRegex := regexp.MustCompile(`href="\/url\?q=(.*?)&amp;`)
		matches := urlRegex.FindAllStringSubmatch(string(body), -1)

		// Print out list of matched URLs
		for _, match := range matches {
			fmt.Println(match[1])
			if _, err = output.WriteString(match[1] + "\n"); err != nil {
				panic(err)
			}
		}

	}

}

func main() {
	banner()
	var help bool
	var wordList string
	var domain string
	var dork bool
	bad_status_code := [10]int{404, 401, 500, 501, 502, 503, 504, 505}

	flag.BoolVar(&help, "h", false, "a breif guide of the tool")
	flag.StringVar(&wordList, "w", "", "relative Path of the wordlist")
	flag.StringVar(&domain, "d", "", "the target domain")
	flag.BoolVar(&dork, "g", false, "use google dorks")

	flag.Parse()

	if help {
		// run guide func
		guide()
	}
	if wordList != "" && domain != "" {
		check_subs(domain, wordList, bad_status_code)
	}
	if domain != "" && dork {
		dorking(domain)
	}
}
