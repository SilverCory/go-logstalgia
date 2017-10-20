package main

import (
	"github.com/SilverCory/go-logstalgia/config"
	"github.com/SilverCory/go-logstalgia/http"
	"strconv"
	"time"
)

func main() {

	conf := &config.LogstalgiaConfig{}
	conf.Load()

	s := http.New(conf)

	go func() {
		for {
			s.Broadcast(&http.LogEntry{
				Time:   strconv.Itoa(int(time.Now().Unix())),
				Path:   "/meme",
				Size:   420,
				IP:     "42.0.42.0",
				Method: "PUT",
				Result: 418,
			})
		}
	}()

	s.Open()
}
