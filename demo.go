package main

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"fmt"
)

type config struct {
	Path string
}

func FindDup(url string,list []string) bool {
	for _,val := range  list {
		if val == url {
			return true
		}
	}
	return false
}

func GetDir(array []string,dulicate []string,url string) ([]string,[]string,error) {
	response,err := http.Get(url)
	if err != nil {
		return array, dulicate, err
	}
	dir := []string{}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()
	lk,err := regexp.Compile(`href="(\d{4}-\d{2}-\d{2}/)"`)
	if err !=nil {
		return array, dulicate, err
	}
	for _,val := range lk.FindAllStringSubmatch(newStr, -1) {
		dir  = append(dir,val[1])
	}
	for _,val := range dir {
		strin := url + val
		fmt.Println(strin)
		if !FindDup(strin,dulicate) {
			array = append(array,strin)
			dulicate = append(dulicate,strin)
		}
	}
	return array,dulicate,nil
}

func GetTxt(array []string,dulicate []string,url string) ([]string,[]string,error){
	response,err := http.Get(url)
	if err != nil {
		return array, dulicate, err
	}
	dir := []string{}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()
	lk,err := regexp.Compile(`href="(.*.all\.txt?)"`)
	if err !=nil {
		return array, dulicate, err
	}
	for _,val := range lk.FindAllStringSubmatch(newStr, -1) {
		dir  = append(dir,val[1])
	}
	for _,val := range dir {
		strin := url + val
		fmt.Println(strin)
		if !FindDup(strin,dulicate) {
			array = append(array,strin)
			dulicate = append(dulicate,strin)
		}
	}
	return array,dulicate,nil
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func CreateFolder(m map[string]string) (error){
	for key,val := range m {
		re1 := regexp.MustCompile("[0-9]+")
		subPath := ""
		for _,val := range re1.FindAllString(key,3) {
			subPath = filepath.Join(subPath, val)
		}
		var userconfig config
		if _, err := toml.DecodeFile("path.toml", &userconfig); err != nil {
			return err
		}
		path := filepath.Join(userconfig.Path, subPath)
		os.MkdirAll(path, os.ModePerm)
		response,err := http.Get(val)
		if err != nil {
			log.Print(err)
		}
		defer response.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		body := buf.Bytes()
		re1 = regexp.MustCompile(`[0-9]+/(.*txt)`)
		subPath = filepath.Join(subPath, val)
		path = filepath.Join(path,re1.FindAllStringSubmatch(val,-1)[0][1])
		fmt.Println(path)
		_ = ioutil.WriteFile(path, body, 0755)

	}
	return nil
}

func main()  {
	defaul := "https://malshare.com/daily/"
	dir := []string{}
	text := []string{}
	dup := []string{}
	dir,dup,err := GetDir(dir,dup,defaul)
	if err != nil {
		log.Println(err)
	}
	text,dup,err = GetTxt(text,dup,defaul)
	if err != nil {
		log.Println(err)
	}
	for {
		fmt.Println("Lengt of dir", len(dir))
		fmt.Println("Lengt of txt", len(text))
		if len(dir) == 0 {
			break
		}
		strin := dir[0]
		dir = remove(dir, 0)
		dir,dup,err = GetDir(dir,dup,strin)
		if err != nil {
			log.Println(err)
		}
		text,dup,err = GetTxt(text,dup,strin)
		if err != nil {
			log.Println(err)
		}
	}
	m := make(map[string]string)
	for _,val := range text {
		fmt.Println(val)
		re := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+/.*")
		if len(re.FindAllString(val,-1)) < 0 {
			m["0000-00-00"] = val
		}
		if len(re.FindAllString(val,-1)) > 0 {
			m[re.FindAllString(val,-1)[0]] = val
		}
	}
	fmt.Println(len(m))
	fmt.Println(m)
	CreateFolder(m)
}
