package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
	"net/http"
)

func main() {
	dbSession, err := r.Connect(r.ConnectOpts{
		Address:  "ec2-54-196-42-226.compute-1.amazonaws.com",
		Database: "shopfeed_dev",
	})
	if err != nil {
		log.Panic(err.Error())
	}

	router := NewRouter(dbSession)

	router.Handle("channel list", channelList)
	router.Handle("channel subscribe messages", channelSubscribeMessages)
	router.Handle("channel unsubscribe messages", channelUnsubscribeMessages)
	router.Handle("channel add message", channelAddMessage)
	router.Handle("vapor add messages", vaporAddMessages)

	http.Handle("/", router)

	fmt.Println("Running FeedServer ...")
	http.ListenAndServe(":4000", nil)
}
