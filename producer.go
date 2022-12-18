package main

import (
	"flag"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeinput := ""
	timeStep := int64(60)
	flag.StringVar(&timeinput, "startfrom", "", "the time start from")
	flag.Parse()
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	from := time.Now().Unix()

	if len(timeinput) > 0 {
		thistime, err := time.Parse("2006-01-02 15:04:05", timeinput)
		if err == nil {
			from = thistime.Unix()
		} else {
			log.Println("time format error")
		}
	}

	log.Println("from:", time.Unix(from, 0).Format("2006-01-02 15:04:05"))

	var to int64
	go func() {
		for {
			timenow := time.Now().Unix()
			if from+timeStep > timenow {
				log.Println("Zzzz...")
				time.Sleep(time.Duration(timeStep) * time.Second)
				continue
			}

			to = from + timeStep - 1

			data := fmt.Sprintf("from:%s->to:%s", time.Unix(from, 0).Format("2006-01-02 15:04:05"), time.Unix(to, 0).Format("2006-01-02 15:04:05"))
			producer.Publish("p1", []byte(data))
			log.Println(data)
			//next loop
			from = to + 1

		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	producer.Stop()

}
