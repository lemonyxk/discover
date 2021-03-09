/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-27 15:40
**/

package main

import (
	"time"

	"github.com/lemoyxk/console"

	"discover/discover"
	"discover/message"
)

func main() {

	var dis = discover.New([]string{"127.0.0.1:11002"})

	dis.Alive("test", "test1").Watch(func(data []*message.ServerInfo) {
		console.Info(data)
	})

	dis.Register("test", "127.0.0.1:1191poo1ii")

	dis.Key("test", "test1").Watch(func(data string) {
		console.Info(data)
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

	select {}
}
