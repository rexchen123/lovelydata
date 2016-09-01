package models

import (
	"github.com/rexchen123/lovelydata/extensions"
	"gopkg.in/mgo.v2"
)

func GetCollection(name string) (*mgo.Session, *mgo.Collection) {
	session := extensions.M.Session.Copy()
	return session, session.DB("").C(name)
}
