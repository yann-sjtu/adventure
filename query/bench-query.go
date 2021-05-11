package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	concurrency []int
	sleepTime   int
	host        string
	test        bool

	proxyTestFetchIp = "http://webapi.http.zhimacangku.com/getip?num=10&type=1&pro=&city=0&yys=0&port=1&time=1&ts=0&ys=0&cs=0&lb=1&sb=0&pb=4&mr=1&regions="
	//proxyFetchIp     = "http://webapi.http.zhimacangku.com/getip?num=400&type=1&pro=&city=0&yys=0&port=1&pack=138919&ts=0&ys=0&cs=0&lb=1&sb=0&pb=4&mr=1&regions=&big_num=2000"
	proxyFetchIp  = "http://webapi.http.zhimacangku.com/getallip?&big_num=1000&type=1&pro=&city=0&yys=0&port=1&pack=138919&ts=0&ys=0&cs=0&lb=1&sb=0&pb=4&mr=3&regions=&username=chukou01&spec=1"
)

func BenchQueryCmd() *cobra.Command {
	// add flags
	cmd := &cobra.Command{
		Use:   "bench-query",
		Short: "",
		Long:  "",
		Run:   benchQuery,
	}
	flags := cmd.Flags()
	flags.IntSliceVarP(&concurrency, "concurrency", "c", []int{1, 1, 1, 1, 1, 1, 1}, "set the number of query concurrent number per second")
	flags.IntVarP(&sleepTime, "sleeptime", "t", 1, "")
	flags.StringVarP(&host, "host", "n", "https://exchaintest.okexcn.com", "")
	flags.BoolVarP(&test, "test", "s", false, "")
	return cmd
}

var startList = []int{0, 0, 0, 0, 0, 0, 0}

func benchQuery(cmd *cobra.Command, args []string) {
	if test {
		proxyFetchIp = proxyTestFetchIp
	}

	if len(concurrency) != 7 {
		panic(fmt.Errorf("concurrent config length should be 7, acutal len: %d", len(concurrency)))
	}
	total := concurrency[0]
	for i := 1; i < 7; i++ {
		startList[i] = startList[i-1] + concurrency[i-1]
		total += concurrency[i]
	}

	ips := QueryProxyIpList()
	for r := 1; ; r++ {
		if r % 15 == 0 {
			newIps := QueryProxyIpList()
			if len(newIps) > 1{
				ips = newIps
			}
		}
		for n := 0; n < 7; n++ {
			reqType := n
			req := generateRequest(reqType)
			for i := 0; i < concurrency[reqType]; i++ {
				go func(round int, num int, typeIndex int) {
					curIndex := startList[typeIndex] + num + round*total
					//fmt.Println(curIndex%len(ips))
					CallWithProxy(req, typeIndex, "http://"+ips[curIndex%len(ips)])
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
				}(r, i, reqType)
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
		req = persistentGasPriceRequest
	case 4:
		req = persistentGetCodeReuqest
	case 5:
		req = EthGetTransactionCount()
	case 6:
		req = EthGetTransactionReceipt()
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
		return nil
	}
	defer resp.Body.Close()
	conent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	urls := string(conent)
	fmt.Println(urls)
	list := strings.Split(urls, "\r\n")
	if len(list) == 1 {
		return nil
	}
	return list[:len(list)-1]
}
