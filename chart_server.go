package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

var yamlConfig *Config
var log *logrus.Logger
var db *sql.DB

type Point struct {
	Symbol string  `json:"symbol"`
	Ltp    float64 `json:"ltp"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   string  `json:"time"`
}

type Points struct {
	Allpoints []Point `json:"points"`
}

func NFOPoints(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.RemoteAddr + " : " + r.URL.String())
	v := r.URL.Query()
	sym := v.Get("symbol")
	start := v.Get("start")
	end := v.Get("end")

  var rows *sql.Rows
	var err error
	var bstart, bend bool
	var startTime, endTime time.Time

	if len(start) == 10 {
		startTime, err = time.Parse("2006-01-02", start)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid start date , please provide yyyy-mm-dd"})
			return
		}
		bstart = true
	}

	if len(end) == 10 {
		endTime, err = time.Parse("2006-01-02", end)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid end date , please provide yyyy-mm-dd"})
			return
		}
    endTime = endTime.Add(time.Hour *  20)
    fmt.Println("20 hours added ", endTime)
		bend = true
	}

	if bstart == false && bend == false {
		rows, err = db.Query("select trad_sym, ltp, open, close, high, low, updated_at from nfo_chart where trad_sym = ? order by updated_at desc", sym)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, no start and end dates"})
			return
		}

	} else if bstart == true && bend == false {
		rows, err = db.Query("select trad_sym, ltp, open, close, high, low, updated_at from nfo_chart where trad_sym = ? and updated_at >= ? order by updated_at desc", sym, startTime)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, no end date"})
			return
		}

	} else if bstart == true && bend == true {
    fmt.Println(startTime)
    fmt.Println(endTime)
		rows, err = db.Query("select trad_sym, ltp, open, close, high, low, updated_at from nfo_chart where trad_sym = ? and updated_at >= ? and updated_at <= ? order by updated_at desc", sym, startTime, endTime)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, start , end dates mentioned"})
			return
		}
	} else {
		log.Error("no start date, only end is mentioned")
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, start date not found"})
		return
	}
	defer rows.Close()

	//	rows, err := db.Query("select trad_sym, ltp, open, close, high, low, updated_at from nfo_chart where trad_sym = ? order by updated_at desc", sym)
	//	if err != nil {
	//		fmt.Println(err)
	//		log.Error(err)
	//		return
	//	}
	//	defer rows.Close()

	var points []Point
	for rows.Next() {
		var open, close, ltp, high, low float64
		var ltt time.Time
		var symbol string

		err := rows.Scan(&symbol, &ltp, &open, &close, &high, &low, &ltt)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "error in query"})
			return
		}
		var pt Point
		pt.Symbol = symbol
		pt.Ltp = ltp
		pt.Open = open
		pt.Close = close
		pt.High = high
		pt.Low = close
		pt.Time = time.Time(ltt).String()
		points = append(points, pt)
	}
	if len(points) == 0 {
		empty := []Point{}
		json.NewEncoder(w).Encode(empty)
		return
	}
	json.NewEncoder(w).Encode(points)

}


func NSEPoints(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.RemoteAddr + " : " + r.URL.String())
	v := r.URL.Query()
	sym := v.Get("symbol")
	start := v.Get("start")
	end := v.Get("end")

  var rows *sql.Rows
	var err error
	var bstart, bend bool
	var startTime, endTime time.Time

	if len(start) == 10 {
		startTime, err = time.Parse("2006-01-02", start)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid start date , please provide yyyy-mm-dd"})
			return
		}
		bstart = true
	}

	if len(end) == 10 {
		endTime, err = time.Parse("2006-01-02", end)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid end date , please provide yyyy-mm-dd"})
			return
		}
    endTime = endTime.Add(time.Hour *  20)
    fmt.Println("20 hours added ", endTime)
		bend = true
	}

	if bstart == false && bend == false {
		rows, err = db.Query("select trad_sym, ltp, open, close, high, low, updated_at from nse_chart where trad_sym = ? order by updated_at desc", sym)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, no start and end dates"})
			return
		}

	} else if bstart == true && bend == false {
		rows, err = db.Query("select trad_sym, ltp, open, close, high, low, updated_at from nse_chart where trad_sym = ? and updated_at >= ? order by updated_at desc", sym, startTime)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, no end date"})
			return
		}

	} else if bstart == true && bend == true {
    fmt.Println(startTime)
    fmt.Println(endTime)
		rows, err = db.Query("select trad_sym, ltp, open, close, high, low, updated_at from nse_chart where trad_sym = ? and updated_at >= ? and updated_at <= ? order by updated_at desc", sym, startTime, endTime)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, start , end dates mentioned"})
			return
		}
	} else {
		log.Error("no start date, only end is mentioned")
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid query, start date not found"})
		return
	}
	defer rows.Close()

	var points []Point
	for rows.Next() {
		var open, close, ltp, high, low float64
		var ltt time.Time
		var symbol string

		err := rows.Scan(&symbol, &ltp, &open, &close, &high, &low, &ltt)
		if err != nil {
			fmt.Println(err)
			log.Error(err)
			json.NewEncoder(w).Encode(map[string]string{"message": "error in query"})
			return
		}
		var pt Point
		pt.Symbol = symbol
		pt.Ltp = ltp
		pt.Open = open
		pt.Close = close
		pt.High = high
		pt.Low = close
		pt.Time = time.Time(ltt).String()
		points = append(points, pt)
	}
	if len(points) == 0 {
		empty := []Point{}
		json.NewEncoder(w).Encode(empty)
		return
	}
	json.NewEncoder(w).Encode(points)

}

func initDB() {
	var dbaddress = yamlConfig.DB.Username + ":" + yamlConfig.DB.Password + "@tcp(127.0.0.1:3306)/" + yamlConfig.DB.Database + "?parseTime=true"
	var err error
	db, err = sql.Open("mysql", dbaddress)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	//defer db.Close()
	var maxConn = yamlConfig.DB.PoolSize
	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(0)
}

func main() {
	args := os.Args
	yamlConfig = getYamlConfig(args[1])
	log = initLogger(yamlConfig.LogFile)
	log.Info("Chart server started with ", os.Getpid())
	initDB()
	router := mux.NewRouter()
	router.HandleFunc("/NFO", NFOPoints).Methods("GET")
	router.HandleFunc("/NSE", NSEPoints).Methods("GET")
	http.ListenAndServe(":8000", router)
}
