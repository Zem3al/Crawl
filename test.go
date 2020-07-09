package main

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"fmt"
)


func main()  {
	response,err := http.Get("https://malshare.com/daily")
	if err != nil {
		log.Print(err)
	}
	dir := []string{}
	text := []string{}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()
	lk := regexp.MustCompile(`\[DIR\](.*)`)
	for _,val := range lk.FindAllString(newStr,-1) {
		re := regexp.MustCompile("href=\".+?\"")
		dir = append(dir,re.FindAllString(val, -1)...)
	}
	lk = regexp.MustCompile(`\[TXT\](.*)`)
	for _,val := range lk.FindAllString(newStr,-1) {
		re := regexp.MustCompile("href=\".+?\"")
		text = append(text,re.FindAllString(val,-1)...)
	}

	for _,val := range dir {
		re := regexp.MustCompile("\"(.*)\"")
		fmt.Println(re.FindStringSubmatch(val)[1])
	}
}
