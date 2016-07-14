package main

import (
	"net/http"

	"github.com/qiniu/log"
	"gopkg.in/siky/config"
	"gopkg.in/siky/db"
	"gopkg.in/siky/route"
)

func init() {
	cfg, err := config.New("./")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutputLevel(cfg.Logger.Level)
}
func main() {

	session, err := db.InitMgo(config.Get().Database.Mongo.URL)
	if err != nil {
		log.Fatal("err", err)
	}
	defer session.Close()
	mHandlers := route.InitHandlers()
	log.Fatal(http.ListenAndServe(":"+config.Get().Server.Port, mHandlers))
}
