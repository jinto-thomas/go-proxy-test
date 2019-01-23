package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net"
	"os"
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

type Broadcast struct {
	connected bool
	conn      net.Conn
	encoder   *gob.Encoder
}

func proxyServer(list *[]Broadcast) {
	var proxyPort = ":" + yamlConfig.Server.PORT
	ln, err := net.Listen("tcp", proxyPort)
	if err != nil {
		fmt.Println(err)
		log.Error(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			continue
		}

		child := Broadcast{true, conn, gob.NewEncoder(conn)}
		*list = append(*list, child)
		fmt.Println("New client added ", conn)

	}
}

func proxyInputFeed(ch <-chan JsonQuote, connList *[]Broadcast) {

	for {
		select {
		case quote := <-ch:
			//fmt.Println("in channel : ", quote)

			for _, client := range *connList {
				if client.connected == true {
					client.encoder.Encode(quote)
				}
			}
			// case <- time.After(time.Second):
			// 	fmt.Println("time out")
			// default:
			// 	fmt.Println("nothing to do in channel")
		}

	}
}

func proxyFeedClient(wg *sync.WaitGroup, ch chan<- JsonQuote) {
	fmt.Println("starting client....")
	var address = yamlConfig.Client.IP + ":" + yamlConfig.Client.PORT
	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println(err)
		log.Error(err)
		wg.Done()
		return
	}

	fmt.Println(err)

	var quote JsonQuote
	decoder := json.NewDecoder(conn)
	for {
		err = decoder.Decode(&quote)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			if err == io.EOF {
				fmt.Println("Source disconnected..")
				break
			}
		} else {
			ch <- quote
		}
	}
}

func main() {

	args := os.Args
	yamlConfig = getYamlConfig(args[1])
	log = initLogger(yamlConfig.LogFile)
	log.Info("Proxy server started with ", os.Getpid())

	messageChannel := make(chan JsonQuote, 10)
	connList := make([]Broadcast, 10)
	go proxyServer(&connList)
	go proxyInputFeed(messageChannel, &connList)
	var wg sync.WaitGroup
	wg.Add(1)
	go proxyFeedClient(&wg, messageChannel)
	wg.Wait()
}
