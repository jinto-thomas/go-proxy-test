package main

import (
	"encoding/gob"
  "encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
)

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
	encoder   *json.Encoder
  id int
}

func proxyServer(list *[]Broadcast) {
	var proxyPort = ":10000"
	ln, err := net.Listen("tcp", proxyPort)
	if err != nil {
		fmt.Println(err)
		return
	}
  id := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
    id += 1
		child := Broadcast{true, conn, json.NewEncoder(conn), id}
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
	var address = "ec2-52-66-242-151.ap-south-1.compute.amazonaws.com" + ":" + "5000"
	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	fmt.Println(err)

	var quote JsonQuote
	//decoder := json.NewDecoder(conn)
  decoder := gob.NewDecoder(conn)
	for {
		err = decoder.Decode(&quote)
		if err != nil {
			fmt.Println(err)
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

	messageChannel := make(chan JsonQuote, 10)
	connList := make([]Broadcast, 10)
	go proxyServer(&connList)
	go proxyInputFeed(messageChannel, &connList)
	var wg sync.WaitGroup
	wg.Add(1)
	go proxyFeedClient(&wg, messageChannel)
	wg.Wait()
}
