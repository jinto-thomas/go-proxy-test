package main

import "encoding/json"
import "fmt"
import "io"
import "net"
import unsafe "unsafe"
import "sync"

//  "encoding/gob"

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

func server() {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleServerConnection(c)
	}
}

func proxyFeedClient(wg *sync.WaitGroup) {
	fmt.Println("starting client....")
	conn, err := net.Dial("tcp", "localhost:3001")

	if err != nil {
		fmt.Println(err)
		fmt.Println("Error occured while trying to connect to the server")
		wg.Done()
		return
	}

	fmt.Println(err)

	buf := make([]byte, 1024*2)
	var packet_len int = 0

	for {
		temp := make([]byte, 256)
		n, erro := conn.Read(temp)
		fmt.Println("bytes read : ", n)
		if erro != nil {
			if erro != io.EOF {
				fmt.Println("Some rror ", erro)
			} else {
				fmt.Println("EOF REACHED.... Source disconnected")
				wg.Done()
				return
			}
		}

		packet_len += n - 1
		copy(buf, temp[:n])

		if temp[n-1] == 10 {
			//got \n
			var quote JsonQuote
			fmt.Println("taking...", packet_len)
			err = json.Unmarshal(buf[:packet_len], &quote)
			if err != nil {
				fmt.Println("json unmarshall error ", err)
			}
			fmt.Println(quote)
			packet_len = 0

			clearbuffer(buf)
		} else {
			fmt.Println("partial data..", packet_len)

			partial := handlePartialRead(buf)
			clearbuffer(buf)
			if partial != nil {

				copy(buf, partial)
				packet_len = len(partial)
			} else {
				packet_len = 0
			}
		}
	}
}

func clearbuffer(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}

func handlePartialRead(buf []byte) []byte {
	var quote JsonQuote
	var start int = 0
	for i, val := range buf {
		if val == 10 {
			err := json.Unmarshal(buf[start:i-1], &quote)
      if err != nil {
        fmt.Println("in partial handler : ", err)
      }
			fmt.Println("from partial read : ", quote)
			start = i
		}
	}

	// remaining data
	if start != len(buf) {
		return buf[start:]
	}

	return nil
}

func handleServerConnection(c net.Conn) {
	// var msg string
	// err := gob.NewDecoder(c).Decode(&msg)
	// if err != nil {
	//   fmt.Println(err)
	// } else {
	//   fmt.Println("Got : ", msg)
	// }

	//  buf := make([]byte, 1024*2)

	temp := make([]byte, 174)

	for {
		fmt.Println("tryig to read...")
		n, err := c.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:: ", err)
				break
			}
		}
		fmt.Println("read bytes ", n, temp)
		//fmt.Println(temp)

		var quote JsonQuote

		error := json.Unmarshal(temp, &quote)
		if error != nil {

		}
		fmt.Println(quote)

		break

	}

	c.Close()
}

func main() {
	var quote JsonQuote
	fmt.Println(unsafe.Sizeof(quote))

	var wg sync.WaitGroup
	wg.Add(1)
	go proxyFeedClient(&wg)
	wg.Wait()

	//go server()

	//var input string
	//fmt.Scanln(&input)
}
