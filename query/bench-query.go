package query

import (
	"fmt"
	"io/ioutil"
	"log"
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
			index := n
			req := generateRequest(index)
			for i := 0; i < concurrency[index]; i++ {
				go func(num int) {
					res, err := CallWithProxy(req, "http://"+ips[rand.Intn(len(ips))])
					if err != nil {
						log.Println("query failed:", err)
						return
					}
					if res.Error != nil {
						log.Println("query result is wrong:", res.Error)
					} else {
						log.Println("query success:", string(res.Result))
					}
				}(i)
			}
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func generateRequest(index int) Request {
	switch index {
	case 0:
		return persistentBlockNumberRequest
	case 1:
		return EthGetBalance()
	case 2:
		return EthGetBlockByNumber()
	case 3:
		return persistentGasPriceRequest
	case 4:
		return persistentGetCodeReuqest
	case 5:
		return EthGetTransactionCount()
	case 6:
		return EthGetTransactionReceipt()
	default:
	}
	return Request{}
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
	list := strings.Split(urls, "\r\n")
	return list[:2000]
}