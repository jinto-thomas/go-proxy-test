package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var yamlConfig *Config
var log *logrus.Logger

type JsonQuote struct {
	Symbol         string  `json:"sym"`
	TradeSymbol    string  `json:"tradSym"`
	Exchange       string  `json:"exc"`
	Ltp            float64 `json:"ltp"`
	Open           float64 `json:"open"`
	Close          float64 `json:"close"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Time           int64   `json:"time"`
	Change         float64 `json:"chg"`
	ChangePer      float64 `json:"chgPer"`
	Ask            float64 `json:"ask"`
	Bid            float64 `json:"bid"`
	BidSize        int64   `json:"askQty"`
	AskSize        int64   `json:"bidQty"`
	OpenInterest   float64 `json:"oi"`
	TotalQty       int64   `json:"tq"`
	InstrumentType string  `json:"type"`
}

func listenToProxy(wg *sync.WaitGroup, chE chan<- JsonQuote, chF chan<- JsonQuote) {

	var address = yamlConfig.Server.IP + ":" + yamlConfig.Server.PORT

	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println(err)
		log.Error(err)
		wg.Done()
		return
	}

	var quote JsonQuote
	decoder := gob.NewDecoder(conn)
	var exc string

	for {
		err = decoder.Decode(&quote)
		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				fmt.Println("Server closed, closing child too")
				log.Error("Server disconnected")
				wg.Done()
				conn.Close()
				break
			}
		} else {
			exc = quote.Exchange
			//fmt.Println(quote)
			if strings.Compare("NSE", exc) == 0 {
				chE <- quote
			} else {
				chF <- quote
			}
			//ch <- quote
		}

	}

}

func updateQuoteDbNSE(ch <-chan JsonQuote) {
	buf := make([]JsonQuote, 200)
	mutex := new(sync.Mutex)

	var dbaddress = yamlConfig.DB.Username + ":" + yamlConfig.DB.Password + "@tcp(127.0.0.1:3306)/" + yamlConfig.DB.Database
	db, err := sql.Open("mysql", dbaddress)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer db.Close()
	var maxConn = yamlConfig.DB.PoolSize
	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(0)

	i := 0
	for {
		select {
		case quote := <-ch:
			buf = append(buf, quote)
			if i > 200 {
				mutex.Lock()
				temp := buf[:i]
				buf = nil
				i = 0
				mutex.Unlock()
				updateValues(&temp, db, "NSE")
				temp = nil
			} else {
				i += 1
			}
			// default:
			// 	fmt.Println("nothing to do in channel")
		}
	}
}

func updateQuoteDbNFO(ch <-chan JsonQuote) {
	buf := make([]JsonQuote, 200)
	mutex := new(sync.Mutex)

	var dbaddress = yamlConfig.DB.Username + ":" + yamlConfig.DB.Password + "@tcp(127.0.0.1:3306)/" + yamlConfig.DB.Database
	db, err := sql.Open("mysql", dbaddress)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer db.Close()
	var maxConn = yamlConfig.DB.PoolSize
	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(0)

	i := 0
	for {
		select {
		case quote := <-ch:
			buf = append(buf, quote)
			if i > 200 {
				mutex.Lock()
				temp := buf[:i]
				buf = nil
				i = 0
				mutex.Unlock()
				updateValues(&temp, db, "NFO")
				temp = nil
			} else {
				i += 1
			}
			// default:
			// 	fmt.Println("nothing to do in channel")
		}
	}
}

func updateValues(buf *[]JsonQuote, db *sql.DB, exc string) {
	var query1 string
	if strings.Compare(exc, "NSE") == 0 {
		query1 = `insert into nse_quote (trad_sym, symbol, ltp, open, close, high, low, changed, changePer, volume, ask, bid, asksize, bidsize, ltt, instrument_type) values `
	} else {
		query1 = `insert into nfo_quote (trad_sym, symbol, ltp, open, close, high, low, changed, changePer, volume, ask, bid, asksize, bidsize, ltt, instrument_type) values `
	}

	tail := ` on duplicate key update ltp = VALUES(ltp), open = VALUES(open), close = VALUES(close), high = VALUES(high), low = VALUES(low),
	changed = VALUES(changed), changePer = VALUES(changePer), volume = VALUES(volume), ask = VALUES(ask), bid = VALUES(bid),
	asksize = VALUES(asksize), bidsize = VALUES(bidsize), ltt = VALUES(ltt), instrument_type = VALUES(instrument_type)`

	var query string
	vals := []interface{}{}

	query += query1
	for _, q := range *buf {
		if q.Ltp == 0 {
			return
		}
		query += "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,from_unixtime(?),?)"
		vals = append(vals, q.TradeSymbol, q.Symbol, q.Ltp, q.Open, q.Close, q.High, q.Low, q.Change, q.ChangePer, q.TotalQty, q.Ask, q.Bid, q.AskSize, q.BidSize, q.Time, q.InstrumentType)
		query += ","
	}

	query = strings.TrimSuffix(query, ",")
	query += tail
	//fmt.Println(query)
	//fmt.Println(vals)

	pstmt, err := db.Prepare(query)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	_, err = pstmt.Exec(vals...)
	if err != nil {
		fmt.Println(err)
		log.Error(err)
	}
	//fmt.Println("updated ", len(*buf))

}

func main() {
	args := os.Args
	yamlConfig = getYamlConfig(args[1])
	log = initLogger(yamlConfig.LogFile)
	log.Info("Quote updater started with ", os.Getpid())
	nsequoteQueueChannel := make(chan JsonQuote, 10)
	nfoquoteQueueChannel := make(chan JsonQuote, 10)

	var wg sync.WaitGroup
	wg.Add(1)
	go updateQuoteDbNSE(nsequoteQueueChannel)
	go updateQuoteDbNFO(nfoquoteQueueChannel)
	go listenToProxy(&wg, nsequoteQueueChannel, nfoquoteQueueChannel)
	wg.Wait()
}
