package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
	"os"
	"github.com/Sirupsen/logrus"
)

var yamlFile string
var log *logrus.Logger

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
	AskSize      int64   `json:"askQty"`
	BidSize      int64   `json:"bidQty"`
	TotalQty     int64   `json:"tq"`
	OpenInterest float64 `json:"oi"`
}

func server(wg *sync.WaitGroup) {
	yamlConfig := getYamlConfig(yamlFile)
	log = initLogger(yamlConfig.LogFile)
	log.Debug("dummy feed started..")
	var port = ":" + yamlConfig.Server.PORT

	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		log.Error(err)
		wg.Done()
		return
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			continue
		}

		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {

	quote1 := JsonQuote{"INFY", "INFY18SEPFUT", "NFO", 714.4, 718.3, 715.9, 731.0, 701.7, 1538890085, 13, 2.1, 715.6, 714.55, 1000, 897, 78263, 8000}
	quote2 := JsonQuote{"BANKBARODA", "BANKBARODA18SEPFUT", "NFO", 712.4, 718.3, 715.9, 731.0, 701.7, 1538890086, 13.4, 2.4, 715.5, 714.40, 980, 997, 74563, 8200}
	quote3 := JsonQuote{"NIFTY", "NIFTY18SEPFUT", "NFO", 12716.4, 11718.3, 11715.9, 11731.0, 11701.7, 1538890087, 13.1, 10.9, 12715.4, 12714.35, 10220, 9230, 752263, 18500}
	quote4 := JsonQuote{"TVSMOTOR", "TVSMOTOR18NOVFUT", "NFO", 507.25, 493.15, 479.85, 519.9, 493.15, 1539164142, 27.399999999999977, 5.710117745128681, 513.15, 511.35, 2000, 2000, 220000, 101000}
	quote5 := JsonQuote{"ICICIBANK","ICICIBANK18DECFUT","NFO",324,313.6,311.05,326,313.6,1539164141,12.949999999999989,4.163317794566787,323.95,323.05,2750,5500,140250,134750}
	quote6 := JsonQuote{"L&TFH","L&TFH18DECFUT","NFO",137,126.25,125,137,126.25,1539164142,12,9.6,0,135.1,0,9000,36000,31500}
	quote7 := JsonQuote{"CEATLTD","CEATLTD18NOVFUT","NFO",1105,1080,1062,1109.55,1079.5,1539164142,43,4.048964218455744,1106.3,1100.05,1050,350,26950,9800}
	quote8 := JsonQuote{"TVSMOTOR","TVSMOTOR18NOVFUT","NFO",507.25,493.15,479.85,519.9,493.15,1539164142,27.399999999999977,5.710117745128681,513.15,511.35,2000,2000,220000,101000}
	quote9 := JsonQuote{"OIL","OIL18NOVFUT","NFO",200,195,191,200.2,191.4,1539164141,9,4.712041884816754,199,197.55,10197,10197,322905,353496}
	quote10 := JsonQuote{"DHFL","DHFL18NOVFUT","NFO",280,253.55,241.75,290,252.3,1539164142,38.25,15.822130299896585,280.85,279.1,3000,3000,654000,1053000}

	list := make([]JsonQuote, 10	)
	list[0] = quote1
	list[1] = quote2
	list[2] = quote3
	list[3] = quote4
	list[4] = quote5
	list[5] = quote6
	list[6] = quote7
	list[7] = quote8
	list[8] = quote9
	list[9] = quote10


	encoder := json.NewEncoder(c)
	fmt.Println(time.Now())
	for {
		for _, quote := range list {
			err := encoder.Encode(quote)
			if err != nil {
				fmt.Println(err)
				log.Error(err)
				return
			}
			time.Sleep(time.Millisecond * 1)
		}
	}
	//c.Close()
}

func main() {
	args := os.Args
	yamlFile = args[1]

	var wg sync.WaitGroup
	wg.Add(1)
	go server(&wg)
	wg.Wait()
}
