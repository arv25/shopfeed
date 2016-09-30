package main

import (
	r "github.com/dancannon/gorethink"
	"log"
	"net/http"
)

type ChannelSubMsgs struct {
	StoreId   string `json:"storeId" gorethink:"store_id,omitempty"`
	ChannelId string `json:"channelId" gorethink:"channel_id,omitempty"`
}

type ChannelAddMsg struct {
	StoreId   string `json:"storeId" gorethink:"store_id,omitempty"`
	ChannelId string `json:"channelId" gorethink:"channel_id,omitempty"`
	Message   string `json:"message" gorethink:"message,omitempty"`
}

func main() {
	dbSession, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "rtsupport",
	})
	if err != nil {
		log.Panic(err.Error())
	}

	router := NewRouter(dbSession)

	router.Handle("channel list", channelList)
	router.Handle("channel subscribe messages", channelSubscribeMessages)
	router.Handle("channel unsubscribe messages", channelUnsubscribeMessages)
	router.Handle("channel add message", channelAddMessage)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
