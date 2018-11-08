package main

import "fmt"

import "math"

import "os"
import "github.com/Sirupsen/logrus"
import "strings"
//import "time"

type F1 struct {
    Age int
    Name string
}

type Formatter2 struct {
  TimestampFormat string
  LogFormat string
}


func (f *Formatter2) Format(entry *logrus.Entry) ([]byte, error) {
  op := f.LogFormat
  tsf := f.TimestampFormat

  op = strings.Replace(op, "%time%", entry.Time.Format(tsf), 1)
  op = strings.Replace(op, "%msg%", entry.Message, 1)
  level := strings.ToUpper(entry.Level.String())
  op = strings.Replace(op, "%lvl%", level,1)
  return []byte(op), nil

  // arr ,_ :=  f.TextFormatter.Format(entry)
  // return append([]byte("☘ "), arr...), nil
}

func main() {
    args := os.Args
    fmt.Println(args[1:])
    yamlConfig := getYamlConfig(args[1])
    var filename = yamlConfig.LogFile

    fmt.Println(yamlConfig)
    fmt.Println("----------")
    fmt.Println(filename)

    f, err := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)

    logger := &logrus.Logger{
      Level: logrus.DebugLevel,
      Formatter: &Formatter2 {
        TimestampFormat: "2006-01-02 15:04:05",
        LogFormat: "%time% - [%lvl%] - %msg%\n",
      },
    }



    //Formatter := new(Formatter2)
    //Formatter.TimestampFormat = "2006-01-02 15:04:05"
    //Formatter.LogFormat =  "%time% - [%lvl%] - %msg%\n"


    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("log file created")
    }
//    logger.SetFormatter(Formatter)
    logger.SetOutput(f)
    logger.Info("log info")
    logger.Debug("hahaha test debug golang log")
    fmt.Println("----------")

}

func Round(x, unit float64) float64 {
    return math.Round(x/unit) * unit
}
