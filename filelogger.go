package main

import (
  "os"
  "strings"
  "github.com/Sirupsen/logrus"
  "sync"
  "fmt"
)

var onceLogger sync.Once
var logger *logrus.Logger

type Formatter struct {
  TimestampFormat string
  LogFormat string
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

  output := f.LogFormat
  tsf := f.TimestampFormat

  output = strings.Replace(output, "%time%", entry.Time.Format(tsf), 1)
  output = strings.Replace(output, "%msg%", entry.Message, 1)
  level := strings.ToUpper(entry.Level.String())
  output = strings.Replace(output, "%lvl%", level,1)
  return []byte(output), nil
}

func initLogger(filename string) *logrus.Logger {
    onceLogger.Do(func() {

      f, err := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
      logger = &logrus.Logger{
        Level: logrus.DebugLevel,
        Formatter: &Formatter {
          TimestampFormat: "2006-01-02 15:04:05",
          LogFormat: "%time% - [%lvl%] - %msg%\n",
        },
      }

      if err != nil {
        fmt.Println(err)
      } else {
        fmt.Println("log file created")
      }
      logger.SetOutput(f)

    })
    return logger
}
