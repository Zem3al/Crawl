package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"regexp"
)

type config struct {
	Default string
	MongoDBHost string
	MongoDBUser string
	MongoDBPwd  string
	Database    string
}

var (
	IsDrop = true
	Config config
)

type Data struct {
	Date string
	md5 string
	sha1 string
	sha256 string
	datasha1 []string
	datasha256 []string
	datamd5 []string
}

type Malware struct {
	//ID    bson.ObjectId `bson:"_id,omitempty"`
	Year  string
	Month string
	Day   string
	Data  string
}


type file struct {
	key     string
	value   string
	session *mgo.Session
}

func Loadconfig() error{
	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		return err
	}
	return nil
}

func GetDir(url string) ([]string, error) {
	array := []string{}
		response, err := http.Get(url)
	if err != nil {
		return array, err
	}
	dir := []string{}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()
	lk, err := regexp.Compile(`href="(\d{4}-\d{2}-\d{2}/)"`)
	if err != nil {
		return array, err
	}
	for _, val := range lk.FindAllStringSubmatch(newStr, -1) {
		dir = append(dir, val[1])
	}
	for _, val := range dir {
		strin := url + val
		fmt.Println(strin)
		array = append(array, strin)
		}

	return array, nil
}

func GetTxt(url string) (Data, error) {
	array := Data{}
	response, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		return array, err
	}
	dir := []string{}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.String()
	lk, err := regexp.Compile(`href="(.*.[^all].txt)"`)
	if err != nil {
		return array, err
	}
	for _, val := range lk.FindAllStringSubmatch(newStr, 3) {
		dir = append(dir, val[1])
	}
	if len(dir) > 0 {
		array.sha1 = fmt.Sprint(url, dir[0])
		array.sha256 = fmt.Sprint(url, dir[1])
		array.md5 = fmt.Sprint(url, dir[2])
	}
	lk, err = regexp.Compile(`(\d{4}-\d{2}-\d{2})`)
	date := lk.FindString(url)
	array.Date = date
	return array, nil
}

func remove(s []string, i int) []string {
	if i >= len(s) || i < 0 {
		return s
	}
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func ConnectDB() (m *mgo.Session,err error) {
	session, _ := mgo.Dial(Config.MongoDBHost)
	session.SetMode(mgo.Monotonic, true)
	if IsDrop {
		err := session.DB(Config.Database).DropDatabase()
		if err != nil {
			return nil,err
		}
	}
	c := session.DB(Config.Database).C("malware")
	index := mgo.Index{
		Key:        []string{"year", "month", "day"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return nil,err
	}
	return session,nil
}

func CreateFolder(key string, val string, session *mgo.Session) error {
	re1 := regexp.MustCompile("[0-9]+")
	lmao := Malware{}
	lmao.Year = re1.FindAllString(key, 3)[0]
	lmao.Month = re1.FindAllString(key, 3)[1]
	lmao.Day = re1.FindAllString(key, 3)[2]
	response, err := http.Get(val)
	if err != nil {
		log.Print(err)
	}
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	body := buf.String()
	lmao.Data = body
	if err != nil {
		panic(err)
	}
	c := session.DB("crawl").C("malware")
	err = c.Insert(&lmao)
	if err != nil {
		return err
	}
	return nil
}

func Crawl(jobs <-chan string, results chan<- Data) {
	for j := range jobs {
		text, _ := GetTxt(j)
		results <- text
	}
}

func WriteFile(job <-chan file, result chan<- string) {
	for j := range job {
		CreateFolder(j.key, j.value, j.session)
		result <- "Done"
	}
}

func RunThreadCrawl(dir []string) ([]Data){
	text := []Data{}
	jobs := make(chan string, len(dir))
	results := make(chan Data)
	for w := 1; w <= 100; w++ {
		go Crawl(jobs, results)
	}
	for _, val := range dir {
		jobs <- val
	}
	close(jobs)
	for a := 1; a <= len(dir); a++ {
		r := <-results
		text = append(text,r)
	}
	fmt.Println(text)
	return text
}

Get data

func main() {
	Loadconfig()
	dir, err := GetDir(Config.Default)
	if err != nil {
		log.Println(err)
		return
	}
	text := RunThreadCrawl(dir)
	fmt.Sprintln(text)
	//fmt.Println(len(m))
	//fmt.Println(m)
	//queq := []file{}
	//sess,err := ConnectDB()
	//if err!= nil {
	//	log.Println(err)
	//	return
	//}
	//for Key, val := range m {
	//	lmao := file{key: Key, value: val, session: sess}
	//	queq = append(queq, lmao)
	//}
	//job := make(chan file, len(queq))
	//result := make(chan string)
	//for w := 1; w <= 100; w++ {
	//	go WriteFile(job, result)
	//}
	//for _, val := range queq {
	//	job <- val
	//}
	//close(job)
	//for a := 1; a <= len(queq); a++ {
	//	<-result
	//}
	//defer sess.Close()
}
