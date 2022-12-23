package main

import (
	//	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type myMessageHandler struct {
	PoolOne      chan struct{}
	WaitGroupOne *sync.WaitGroup
}

// HandleMessage implements the Handler interface.
func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	// do whatever actual message processing is desired
	//发送一个阻塞到队列
	for {
		if len(h.PoolOne) == 2 {
			//log.Println("chan limit 2")
			m.Touch()
			time.Sleep(time.Second * 10)
		} else {
			log.Println("chan add 1 ")
			h.PoolOne <- struct{}{}
			break
		}

	}
	h.WaitGroupOne.Add(1)
	go func() {
		h.processMessage(m.Body)
		h.WaitGroupOne.Done()
		<-h.PoolOne
	}()
	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}
func (h *myMessageHandler) processMessage(data []byte) error {
	log.Println("message incoming...")
	time.Sleep(120 * time.Second)
	log.Println("message->", string(data))
	return nil
}
func main() {
	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("p1", "channel", config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	myhd := &myMessageHandler{}
	myhd.PoolOne = make(chan struct{}, 2)
	myhd.WaitGroupOne = &sync.WaitGroup{}
	consumer.AddHandler(myhd)

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	//err = consumer.ConnectToNSQLookupd("localhost:4161")
	err = consumer.ConnectToNSQD("localhost:4150")
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	consumer.Stop()
	//myhd.WaitGroupOne.Wait()
	// Gracefully stop the consumer.
}
