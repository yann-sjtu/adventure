package query

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	concurrency int
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
	flags.IntVarP(&concurrency, "concurrency", "c",1, "set the number of query num")
	flags.IntVarP(&sleepTime, "sleeptime", "t",1, "set the number of query num")
	flags.StringVarP(&host, "host", "n","", "set the number of query num")
	return cmd
}

func benchQuery(cmd *cobra.Command, args []string) {
	list := QueryProxyIpList()

	for {
		time.Sleep(time.Duration(sleepTime) * time.Second)

		//EthBlockNumber()
		//EthGetBalance()
		//EthGetBlockByNumber()
		//EthGasPrice()
		//EthGetCode()
		//EthGetTransactionCount()
		//EthGetTransactionReceipt()
		res, err := CallWithProxy(EthBlockNumber(), "http://"+list[0])
		if err != nil {
			log.Println("query failed:", err)
			continue
		}
		if res.Error != nil {
			log.Println("query result is wrong:", res.Error)
		} else {
			log.Println("query success:", string(res.Result))
		}
	}
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