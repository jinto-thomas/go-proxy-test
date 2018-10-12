package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net"
	"strings"
	"sync"
)

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

func listenToProxy(wg *sync.WaitGroup, ch chan<- JsonQuote) {
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
			if err == io.EOF {
				fmt.Println("Server closed, closing child too")
				wg.Done()
				conn.Close()
				break
			}
		} else {
			//fmt.Println(quote)
			ch <- quote
		}

	}

}

func updateQuoteDb(ch <-chan JsonQuote) {
	buf := make([]JsonQuote, 200)
	mutex := new(sync.Mutex)

	db, err := sql.Open("mysql", "root:roadrunner@tcp(127.0.0.1:3306)/quote")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	i := 0
	for {
		select {
		case quote := <-ch:
			//fmt.Println("index ", i)
			buf = append(buf, quote)
			if i > 200 {
				mutex.Lock()
				temp := buf[:i]
				buf = nil
				i = 0
				mutex.Unlock()
				updateValues(&temp, db)
				temp = nil
			} else {
				i += 1
			}
			// default:
			// 	fmt.Println("nothing to do in channel")
		}
	}
}

func updateValues(buf *[]JsonQuote, db *sql.DB) {
	query1 := `insert into nse_quote (trad_sym, symbol, ltp) values `

	tail := " on duplicate key update ltp = VALUES(ltp)"

	var query string
	vals := []interface{}{}

	query += query1
	for _, q := range *buf {
		if q.Ltp == 0 {
			return
		}
		query += "(?,?,?)"
		vals = append(vals, q.TradeSymbol, q.Symbol, q.Ltp)
		query += ","
		//vals = append(vals, q.Ltp, q.Open, q.Close, q.High, q.Low, q.Change, q.ChangePer, q.TotalQty, q.Ask, q.Bid, q.AskSize, q.BidSize, q.TradeSymbol)
	}

	query = strings.TrimSuffix(query, ",")
	query += tail
	//fmt.Println(query)
	//fmt.Println(vals)

	db.Begin()
	pstmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	_, err = pstmt.Exec(vals...)
	fmt.Println("updated ", len(*buf))

}

func main() {
	quoteQueueChannel := make(chan JsonQuote, 10)

	var wg sync.WaitGroup
	wg.Add(1)
	go updateQuoteDb(quoteQueueChannel)
	go listenToProxy(&wg, quoteQueueChannel)
	wg.Wait()
}
