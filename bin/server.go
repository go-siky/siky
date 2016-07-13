package main

import (
	"log"
	"net/http"
	"time"

	"gopkg.in/siky/db"
	"gopkg.in/siky/route"
)

func main() {
	session, err := db.InitMgo("10.3.30.183")
	if err != nil {
		log.Println("err")
	}
	defer session.Close()

	time.Sleep(time.Second * 30)
	mHandlers := route.InitHandlers()
	log.Fatal(http.ListenAndServe(":8080", mHandlers))

}
