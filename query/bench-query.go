package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	concurrency []int
	sleepTime   int
	host        string

	proxyFetchIp = "http://webapi.http.zhimacangku.com/getip?num=400&type=1&pro=&city=0&yys=0&port=1&pack=138919&ts=0&ys=0&cs=0&lb=1&sb=0&pb=4&mr=1&regions=&big_num=2000"
)

func BenchQueryCmd() *cobra.Command {
	// add flags
	cmd := &cobra.Command{
		Use:   "bench-query",
		Short: "",
		Long:  "",
		Run:  benchQuery,
	}
	flags := cmd.Flags()
	flags.IntSliceVarP(&concurrency, "concurrency", "c", []int{1,1,1,1,1,1,1}, "set the number of query concurrent number per second")
	flags.IntVarP(&sleepTime, "sleeptime", "t",1, "set the number of query num")
	flags.StringVarP(&host, "host", "n","", "set the number of query num")
	return cmd
}

func benchQuery(cmd *cobra.Command, args []string) {
	if len(concurrency) != 7 {
		panic(fmt.Errorf("concurrent config length should be 7, acutal len: %d", len(concurrency)))
	}

	ips := QueryProxyIpList()

	for {
		for n := 0; n < 7; n++ {
			reqType := n
			req := generateRequest(reqType)
			for i := 0; i < concurrency[reqType]; i++ {
				go func(num int) {
					CallWithProxy(req, reqType, "http://"+ips[rand.Intn(len(ips))])
					//res, err := CallWithProxy(req, reqType, "http://"+ips[rand.Intn(len(ips))])
					//if err != nil {
					//	log.Println("query failed:", err)
					//	return
					//}
					//if res.Error != nil {
					//	log.Println("query result is wrong:", res.Error)
					//} else {
					//	log.Println("query success:", string(res.Result))
					//}
				}(i)
			}
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func generateRequest(index int) []byte {
	var req Request
	switch index {
	case 0:
		req = persistentBlockNumberRequest
	case 1:
		req = EthGetBalance()
	case 2:
		req = EthGetBlockByNumber()
	case 3:
		req =  persistentGasPriceRequest
	case 4:
		req =  persistentGetCodeReuqest
	case 5:
		req =  EthGetTransactionCount()
	case 6:
		req =  EthGetTransactionReceipt()
	default:
		req = persistentBlockNumberRequest
	}

	postBody, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return postBody
}

func QueryProxyIpList() []string {
	resp, err := http.Get(proxyFetchIp)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	conent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	urls := string(conent)
	fmt.Println(urls)
	list := strings.Split(urls, "\r\n")
	return list[:2000]
}