package types

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Ticker struct {
	InstrumentId string `json:"instrument_id"`
	Last         string `json:"last"`
	//BestBid        string `json:"best_bid"`
	//BestAsk        string `json:"best_ask"`
	//Open24h        string `json:"open_24h"`
	//High24h        string `json:"high_24h"`
	//Low24h         string `json:"low_24h"`
	//BaseVolume24h  string `json:"base_volume_24h"`
	//QuoteVolume24h string `json:"quote_volume_24h"`
	//Timestamp      string `json:"timestamp"`
}

//"https://www.okex.com/api/spot/v3/instruments/OKB-BTC/ticker"
func QueryOneTickerPrice(instrumentId string) float64 {
	response, err := http.Get(okexUrl + "/api/spot/v3/instruments/" + instrumentId + "/ticker")
	if err != nil {
		log.Println("query ticker failed:", err.Error())
		return 0.0
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("read ticker failed:", err.Error())
		return 0.0
	}

	var ticker Ticker
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		log.Println("unmarshal ticker json failed:", err.Error(), ". body:", string(body))
		return 0.0
	}

	price, _ := strconv.ParseFloat(ticker.Last, 64)
	if price == 0.0 {
		log.Println("[warning] price from ticker is zero")
	}
	return price
}

type OrderMsg struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	DetailMsg string `json:"detail_msg"`
	Data      struct {
		Data      []OrderData `json:"data"`
		ParamPage ParamPage   `json:"param_page"`
	} `json:"data"`
}

type OrderData struct {
	Txhash         string `json:"txhash"`
	OrderID        string `json:"order_id"`
	Sender         string `json:"sender"`
	Product        string `json:"product"`
	Side           string `json:"side"`
	Price          string `json:"price"`
	Quantity       string `json:"quantity"`
	Status         int    `json:"status"`
	FilledAvgPrice string `json:"filled_avg_price"`
	RemainQuantity string `json:"remain_quantity"`
	Timestamp      int    `json:"timestamp"`
	IsFind         bool
}

type ParamPage struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

//https://www.okex.com/okexchain/v1/order/list/open?product="okt_tusdk&page=1&per_page=300&address=okchain1mdzw4um5qpk2sc3vay2jrq0xx52crr0ydsks5m
func QueryOrders(product string, address string) *OrderMsg {
	req, err := http.NewRequest("GET", okexUrl+"/okexchain/v1/order/list/open", nil)
	if err != nil {
		log.Println("create req msg failed:", err.Error())
		return nil
	}

	q := req.URL.Query()
	q.Add("product", product)
	q.Add("page", "1")
	q.Add("per_page", "300")
	q.Add("address", address)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("query order msg failed:", err.Error())
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read order msg failed:", err.Error())
		return nil
	}

	var msg OrderMsg
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Println("unmarshal order msg json failed:", err.Error(), ". body:", body)
		return nil
	}
	return &msg
}
