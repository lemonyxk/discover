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
)

func main() {

	var dis = discover.New([]string{"127.0.0.1:11002"})

	dis.Alive("test", "test1").Watch(func(data []*message.ServerInfo) {
		for i := 0; i < len(data); i++ {
			console.Info("server:", data[i])
		}
	})

	dis.Register("test", "127.0.0.1:1191poo1ii")

	dis.Key("test", "test1").Watch(func(op message.Op) {
		console.Infof("%+v\n", op)
	})

	time.AfterFunc(time.Second, func() {
		console.Info(dis.Delete("test"))
	})

	time.AfterFunc(2*time.Second, func() {
		console.Info(dis.Set("test", "set test"))
	})

	time.AfterFunc(time.Second, func() {
		console.Info(dis.Delete("test1"))
	})

	time.AfterFunc(2*time.Second, func() {
		console.Info(dis.Set("test1", "set test1"))
	})

	time.AfterFunc(3*time.Second, func() {
		console.Info(dis.Get("test1"))
	})

	select {}
}
