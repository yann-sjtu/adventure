package query

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	concurrency []int
	sleepTime   int
	host        string
	test        bool

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
	flags.IntSliceVarP(&concurrency, "concurrency", "g", []int{1, 1, 1, 1, 1, 1, 1}, "set the number of query concurrent number per second")
	flags.IntVarP(&sleepTime, "sleeptime", "t", 1, "")
	flags.StringVarP(&host, "host", "o", "https://exchaintestrpc.okex.org", "")
	flags.BoolVarP(&test, "test", "s", false, "")
	return cmd
}

var startList = []int{0, 0, 0, 0, 0, 0, 0}

func benchQuery(cmd *cobra.Command, args []string) {

	if len(concurrency) != 7 {
		panic(fmt.Errorf("concurrent config length should be 7, acutal len: %d", len(concurrency)))
	}
	total := concurrency[0]
	for i := 1; i < 7; i++ {
		startList[i] = startList[i-1] + concurrency[i-1]
		total += concurrency[i]
	}

	for r := 1; ; r++ {
		for n := 0; n < 7; n++ {
			reqType := n
			for i := 0; i < concurrency[reqType]; i++ {
				go func(round int, num int, typeIndex int) {
					req := generateRequest(reqType)
					//fmt.Println(curIndex%len(ips))
					CallWithProxy(req, typeIndex, "")
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

//func QueryProxyIpList() []string {
//	resp, err := http.Get(proxyFetchIp)
//	if err != nil {
//		return nil
//	}
//	defer resp.Body.Close()
//	conent, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//	urls := string(conent)
//	fmt.Println(urls)
//	list := strings.Split(urls, "\r\n")
//	if len(list) == 1 {
//		return nil
//	}
//	return list[:len(list)-1]
//}
