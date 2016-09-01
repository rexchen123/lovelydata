package extensions

import (
	"github.com/adjust/redismq"
	"github.com/astaxie/beego"
)

type ZhihuQueue struct {
	Name     string
	Queue    *redismq.Queue
	Consumer *redismq.Consumer
}

var Q *ZhihuQueue = nil

func initQueue() {
	Q = &ZhihuQueue{
		Name:     "zhihu",
		Queue:    nil,
		Consumer: nil,
	}
	Q.Queue = redismq.CreateQueue("localhost", "6379", "", 9, Q.Name)
	consumer, err := Q.Queue.AddConsumer("consumer")
	if err != nil {
		beego.Error(err)
	} else {
		Q.Consumer = consumer
	}
}
