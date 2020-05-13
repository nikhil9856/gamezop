package main

import (
	"log"
	"os"
	"sync"

	"github.com/gamezop/model"
	"github.com/gamezop/util"
	nsq "github.com/nsqio/go-nsq"
)

func main() {
	var err error
	wg := &sync.WaitGroup{}
	wg.Add(1)
	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer(model.NsqTopic, model.NsqChannel, config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		err = util.PushEventToDatabase(message.Body)
		wg.Done()
		return err
	}))
	if err = q.ConnectToNSQD("127.0.0.1:4150"); err != nil {
		log.Println(model.NsqConnectionError)
		log.Println(err)
		os.Exit(0)
	}
	wg.Wait()
}
