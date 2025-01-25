/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-27 15:40
**/

package main

import (
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/discover"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/store"
)

func main() {

	// c := store.Message{
	// 	Op:  store.Delete,
	// 	Key: "test1",
	// 	Value: []byte("test1 value"),
	// }
	//
	// console.Info(store.Parse(store.Build(c)))

	var dis = discover.New("127.0.0.1:11002").Connect()

	dis.OnClose(func() {
		console.Info.Log("alive close")
	})

	var alive = dis.Alive("test", "test1")

	alive.Watch(func(name string, serverInfo []*message.ServerInfo) {
		for i := 0; i < len(serverInfo); i++ {
			console.Info.Logf("%s %+v", name, serverInfo[i])
		}
	})

	dis.Register(func() message.ServerInfo {
		return message.ServerInfo{
			Name: "test",
			Addr: "127.0.0.1:1191poo1ii",
		}
	})

	var key = dis.Key("test", "test1")

	key.Watch(func(message *store.Message) {
		console.Info.Log(message)
	})

	time.AfterFunc(time.Second, func() {
		var res, err = dis.Delete("test")
		console.Info.Logf("%s %s", res, err)
	})

	time.AfterFunc(2*time.Second, func() {
		var res, err = dis.Set("test1", []byte("set test1"))
		console.Info.Logf("%s %s", res, err)
	})

	time.AfterFunc(time.Second, func() {
		var res, err = dis.Delete("test1")
		console.Info.Logf("%s %s", res, err)
	})

	time.AfterFunc(2*time.Second, func() {
		var res, err = dis.Set("test2", []byte("set test2"))
		console.Info.Logf("%s %s", res, err)
	})

	time.AfterFunc(3*time.Second, func() {
		var res, err = dis.Get("test1")
		console.Info.Logf("%s %s", res, err)
	})

	time.AfterFunc(3*time.Second, func() {
		var res, err = dis.Get("test2")
		console.Info.Logf("%s %s", res, err)
	})

	time.AfterFunc(2*time.Second, func() {
		var all, err = dis.All()
		if err != nil {
			console.Error.Logf("%s", err)
			return
		}

		console.Info.Log(all)
	})

	select {}
}
