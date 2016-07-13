package db

import "gopkg.in/mgo.v2"

//MgoSession is global mongo session
var MgoSession *mgo.Session

//InitMgo for initialize mongo session
func InitMgo(mongoURL string) (session *mgo.Session, err error) {

	session, err = mgo.Dial(mongoURL)
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return
}
