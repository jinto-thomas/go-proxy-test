package main

import "fmt"
import "encoding/gob"
import "net"
import "sync"
import "io"

type JsonQuote struct {
	Symbol       string  `json:"sym"`
	TradeSymbol  string  `json:"tradSym"`
	Exchange     string  `json:"exc"`
	Ltp          float64 `json:"ltp"`
	Open         float64 `json:"open"`
	Close        float64 `json:"close"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Time         int64   `json:"time"`
	Change       float64 `json:"chg"`
	ChangePer    float64 `json:"chgPer"`
	Ask          float64 `json:"ask"`
	Bid          float64 `json:"bid"`
	BidSize      int64   `json:"askQty"`
	AskSize      int64   `json:"bidQty"`
	OpenInterest float64 `json:"oi"`
	TotalQty     int64   `json:"tq"`
}

func listenToProxy(wg *sync.WaitGroup) {
	conn, err := net.Dial("tcp", "localhost:5000")

	if err != nil {
		fmt.Println(err)
		fmt.Println("Error occured while trying to connect to the server")
		wg.Done()
		return
	}

	var quote JsonQuote
	decoder := gob.NewDecoder(conn)

	for {
		err = decoder.Decode(&quote)
		if err != nil {
			fmt.Println(err)
			if (err == io.EOF) {
				fmt.Println("Server closed, closing child too")
				wg.Done()
				conn.Close()
				break
			}
		} else {
			fmt.Println(quote)
		}

	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go listenToProxy(&wg)
	wg.Wait()
}
