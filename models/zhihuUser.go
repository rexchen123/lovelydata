package models

import (
	"fmt"
	"github.com/astaxie/beego"
	. "github.com/rexchen123/lovelydata/extensions"
	. "github.com/rexchen123/lovelydata/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ZhihuUser struct {
	Id               bson.ObjectId `bson:"_id"`
	ZhihuId          string        `bson:"zhihuId"`
	Name             string        `bson:"name"`
	Address          string        `bson:"address"`
	Avatar           string        `bson:"avatar"`
	Business         string        `bson:"business"`
	Gender           string        `bson:"gender"`
	Education        string        `bson:"education"`
	Major            string        `bson:"major"`
	Description      string        `bson:"description"`
	FolloweesCount   int           `bson:"followeesCount"`
	FollowersCount   int           `bson:"followersCount"`
	SpecialCount     int           `bson:"specialCount"`
	FollowTopicCount int           `bson:"followTopicCount"`
	ApprovalCount    int           `bson:"approvalCount"`
	ThankCount       int           `bson:"thankCount"`
	AnswerCount      int           `bson:"answerCount"`
	ArticleCount     int           `bson:"articleCount"`
	PvCount          int           `bson:"pvCount"`
	StartedCount     int           `bson:"startedCount"`
	PublicEditCount  int           `bson:"publicEditCount"`
	AskCount         int           `bson:"askCount"`
}

func WaitToCrawlUser() {
	go func(Q *ZhihuQueue) {
		tick := time.Tick(time.Second * 1)
		for {
			select {
			case <-tick:
				getUser(Q)
			}
		}

	}(Q)
}

func getUser(Q *ZhihuQueue) {
	data, err := Q.Consumer.Get()
	if data != nil {
		beego.Info(data.Payload)
		fmt.Println(data.Payload)
		go RecordUser(data.Payload)
		go GetFollowers(data.Payload)
		go GetFollowees(data.Payload)
		err = data.Ack()
		if err != nil {
			beego.Error(err)
		}
	}
	if err != nil {
		beego.Error(err)
	}
}

func GetFollowers(zhihuId string) {
	GetFollowerList(zhihuId, "followers")
}

func GetFollowees(zhihuId string) {
	GetFollowerList(zhihuId, "followees")
}

func RecordUser(zhihuId string) {
	user := FetchUserInfo(zhihuId)
	zhihuUser := ZhihuUser{
		Id:               bson.NewObjectId(),
		ZhihuId:          user.ZhihuId,
		Name:             user.Name,
		Address:          user.Address,
		Avatar:           user.Avatar,
		Business:         user.Business,
		Gender:           user.Gender,
		Education:        user.Education,
		Major:            user.Major,
		Description:      user.Description,
		FolloweesCount:   user.FolloweesCount,
		FollowersCount:   user.FollowersCount,
		SpecialCount:     user.SpecialCount,
		FollowTopicCount: user.FollowTopicCount,
		ApprovalCount:    user.ApprovalCount,
		ThankCount:       user.ThankCount,
		AnswerCount:      user.AnswerCount,
		ArticleCount:     user.ArticleCount,
		PvCount:          user.PvCount,
		StartedCount:     user.StartedCount,
		PublicEditCount:  user.PublicEditCount,
		AskCount:         user.AskCount,
	}

	session, collection := GetCollection("zhihuUser")
	defer session.Close()

	err := collection.Insert(zhihuUser)
	if err != nil {
		msg := fmt.Sprintf("Error record zhihu user failed, error: %v", err)
		beego.Warn(msg)
	}
}
