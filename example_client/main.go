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
		console.Info("alive close")
	})

	var alive = dis.Alive("test", "test1")

	alive.Watch(func(name string, serverInfo []*message.ServerInfo) {
		for i := 0; i < len(serverInfo); i++ {
			console.Info("server:", serverInfo[i])
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
		console.Infof("%+v\n", message)
	})

	time.AfterFunc(time.Second, func() {
		console.Info(dis.Delete("test"))
	})

	time.AfterFunc(2*time.Second, func() {
		console.Info(dis.Set("test1", []byte("set test1")))
	})

	time.AfterFunc(time.Second, func() {
		console.Info(dis.Delete("test1"))
	})

	time.AfterFunc(2*time.Second, func() {
		console.Info(dis.Set("test2", []byte("set test2")))
	})

	time.AfterFunc(3*time.Second, func() {
		console.Info(dis.Get("test1"))
	})

	time.AfterFunc(3*time.Second, func() {
		console.Info(dis.Get("test2"))
	})

	time.AfterFunc(2*time.Second, func() {
		var all, err = dis.All()
		if err != nil {
			console.Error(err)
			return
		}

		console.Infof("%s\n", all)
	})

	select {}
}
