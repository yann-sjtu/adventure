package scenario

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	URL 		string
	AccountList string
	Account 	string
}

/**
该文件提供些基础的util 函数
 */
func PanicErr(err error){
	PanicErr(err)
}

/**
读取文件中的每一行，放到[]string中
 */
func ReadDataFromFile(path string)(lines []string){
	f, err := os.Open(path)
	PanicErr(err)
	rd := bufio.NewReader(f)
	lines = make([]string, 0)
	num := 0
	for true{
		line, err := rd.ReadString('\n')
		if err != nil || err == io.EOF{
			break
		}
		lines = append(lines,strings.TrimSpace(line))
		num++
	}
	log.Printf("Read lines %d from file", num)
	return
}

/**
读取json配置映射到struct中
 */
func ReadJson2Struct(filepath string, obj interface{}){
	file, ferr := os.Stat(filepath)
	if ferr != nil{
		log.Printf("file path does not exist, %s", filepath)
	}
	log.Printf("%v", file)
	b, err := ioutil.ReadFile(filepath)
	PanicErr(err)
	err1 := json.Unmarshal(b, &obj)
	if err != nil {
		log.Printf("error is: %s", err1)
	}
}


const (
	success = "success"
	fail    = "failed"
)

func DoPost(url string, postBody []byte) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	startTime := time.Now()
	resp, respErr := client.Do(req)
	elapsed := time.Since(startTime)

	if respErr != nil {
		log.Println(postBody, strconv.FormatInt(elapsed.Milliseconds(),10) + "ms", fail, respErr)
		return nil, respErr
	}

	log.Println(postBody, strconv.FormatInt(elapsed.Milliseconds(),10) + "ms", success)
	defer resp.Body.Close()
	return resp, nil
}

func DoGet(url string, body []byte){

}

