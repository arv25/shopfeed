package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
	math "math/rand"
	"strconv"
	"time"
)

func vaporAddMessages(client *Client, data interface{}) {
	var clientData map[string]string
	mapstructure.Decode(data, &clientData)
	limit, err := strconv.Atoi(clientData["count"])
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}
	client.send <- Message{"Hydrating some messages for you every second ...", ""}

	go func() {

		for i := 0; i < limit; i++ {
			msg := make(map[string]string)

			channelIds := []string{"a301004a-824f-48e4-b5fe-e2c15aa2f085", "63110664-d13e-4b1c-ae61-c67f544ecd04", "cb3f5892-d778-44e6-8372-ac47028c29c0", "6b86f679-aa9e-4297-9779-23383522a718", "2f931159-9d98-4455-8e0b-fa44a05e2284", "85b552da-b553-44e1-903d-125a359235df", "53fe5358-ee3a-44c9-893d-77686f8d95a3", "efe2bdb2-f04b-4ca4-b678-a314b921cf2f", "57be3d6f-b614-49f8-a104-e33fa47e0fbc"}
			msg["channelId"] = channelIds[math.Intn(8)]

			storeIds := []string{"ABC1234", "DEF4567"}
			msg["storeId"] = storeIds[math.Intn(1)]

			sources := []string{"Backoffice", "Pocket"}
			msg["source"] = sources[math.Intn(1)]

			times := []string{"9/28/2016 22:15:44 UTC", "9/27/2016 09:29:44 UTC", "9/30/2016 18:30:44 UTC", "9/29/2016 10:13:64 UTC"}
			msg["time"] = times[math.Intn(3)]

			types := []string{"EventTypeA", "EventTypeB"}
			msg["type"] = types[math.Intn(1)]

			users := []string{"Hankster", "Billy", "Jessica", "Donatello", "Luigi", "Matthais", "Victor", "Kyle", "Brian", "Julie", "Bubba"}
			msg["userName"] = users[math.Intn(10)]

			messages := []string{"Some Rand Message", "I love cupcakes", "Lots of returns today", "Something happened that's important."}
			msg["message"] = messages[math.Intn(3)]

			fmt.Printf("%#v\n", msg)

			err := r.Table("messages").
				Insert(msg).
				Exec(client.dbSession)
			if err != nil {
				client.send <- Message{"error", err.Error()}
			}
			time.Sleep(time.Second)
		}

		client.send <- Message{"Hydrated some messages for you.", ""}
	}()
}
