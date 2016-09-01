package extensions

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"gopkg.in/mgo.v2"
	"time"
)

type MgoClient struct {
	conf    map[string]string
	Session *mgo.Session
}

var M *MgoClient = nil

func initMongo(conf config.Configer) {
	M = &MgoClient{
		conf:    make(map[string]string),
		Session: nil,
	}
	M.initWithConf(conf)
}

func (self *MgoClient) initWithConf(conf config.Configer) {
	host := conf.String("mongo::host")
	port := conf.String("mongo::port")
	user := conf.String("mongo::user")
	pwd := conf.String("mongo::pwd")
	db := conf.String("mongo::db")

	mhost := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, pwd, host, port, db)
	beego.Info(mhost)

	mgosession, err := mgo.Dial(mhost)

	if err != nil {
		msg := fmt.Sprintf("Failed to connect MongoDB with host: %s, err: %v", mhost, err)
		beego.Error(msg)
	} else {
		msg := fmt.Sprintf("Successfully connected to mongodb, host %s port %s user %s", host, port, user)
		beego.Info(msg)
	}

	mgosession.SetMode(mgo.Monotonic, true)
	beego.Info("Enabled mongodb Monotonic mode")

	self.Session = mgosession
	self.conf = map[string]string{
		"host": host,
		"port": port,
		"user": user,
		"pwd":  pwd,
		"db":   db,
	}

	go func(self *MgoClient) {
		tick := time.Tick(time.Second * 5)

		for {
			select {
			case <-tick:
				ensureConnected(self, conf)
			}
		}

	}(self)
}

func ensureConnected(client *MgoClient, conf config.Configer) {
	if nil == client.Session {
		client.initWithConf(conf)
	} else {
		err := client.Session.Ping()

		if nil != err {
			msg := fmt.Sprintf("Failed to Ping mongo server, seems the connection is down, %v", err)
			beego.Warn(msg)

			client.Session.Close()
			client.initWithConf(conf)
		}
	}
}
