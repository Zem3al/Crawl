package main

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"fmt"
)

func FindDup(url string,list []string) bool {
	for _,val := range  list {
		if val == url {
			return true
		}
	}
	return false
}

func GetDir(array *[]string,dulicate *[]string,url string) {
	response,err := http.Get(url)
	if err != nil {
		log.Print(err)
	}
	dir := []string{}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()
	lk := regexp.MustCompile(`\[DIR\](.*)`)
	for _,val := range lk.FindAllString(newStr,-1) {
		re := regexp.MustCompile("href=\".+?\"")
		dir = append(dir,re.FindAllString(val, -1)...)
	}
	for _,val := range dir {
		re := regexp.MustCompile("\"(.*)\"")
		strin := url + re.FindStringSubmatch(val)[1]
		if !FindDup(strin,*dulicate) {
			*array = append(*array,strin)
			*dulicate = append(*dulicate,strin)
		}
	}
}

func GetTxt(array *[]string,dulicate *[]string,url string) {
	response,err := http.Get(url)
	if err != nil {
		log.Print(err)
	}
	dir := []string{}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()
	lk := regexp.MustCompile(`\[TXT\](.*)`)
	for _,val := range lk.FindAllString(newStr,-1) {
		re := regexp.MustCompile("href=\".+?\"")
		dir = append(dir,re.FindAllString(val, -1)...)
	}
	for _,val := range dir {
		re := regexp.MustCompile("\"(.*)\"")
		strin := url + re.FindStringSubmatch(val)[1]
		if !FindDup(strin,*dulicate) {
			*array = append(*array,strin)
			*dulicate = append(*dulicate,strin)
		}
	}
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func main()  {
	defaul := "https://malshare.com/daily"
	dir := []string{}
	text := []string{}
	dup := []string{}
	GetDir(&dir,&dup,defaul)
	GetTxt(&text,&dup,defaul)
	i :=1
	for {
		fmt.Println(i)
		i++
		fmt.Println(len(dir))
		if len(dir) == 0 {
			break
		}
		strin := dir[0]
		dir = remove(dir,0)
		GetDir(&dir,&dup,strin)
		GetTxt(&text,&dup,strin)
	}
	fmt.Println(text)
}
