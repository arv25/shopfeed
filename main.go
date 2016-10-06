package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
	"net/http"
)

type ChannelSubMsgs struct {
	StoreId   string `json:"storeId" gorethink:"storeId,omitempty"`
	ChannelId string `json:"channelId" gorethink:"channelId,omitempty"`
}

type ChannelAddMsg struct {
	StoreId   string `json:"storeId" gorethink:"storeId,omitempty"`
	ChannelId string `json:"channelId" gorethink:"channelId,omitempty"`
	Source    string `json:"source" gorethink:"source,omitempty"`
	Time      string `json:"time" gorethink:"time,omitempty"`
	Type      string `json:"type" gorethink:"type,omitempty"`
	UserName  string `json:"userName" gorethink:"userName,omitempty"`
	Message   string `json:"message" gorethink:"message,omitempty"`
}

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
