/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-26 11:46
**/

package client

import (
	"time"

	"github.com/lemoyxk/console"
	client2 "github.com/lemoyxk/kitty/socket/udp/client"

	"discover/app"
)

func Start(host string, fn func()) {

	var client = &client2.Client{Host: host, Reconnect: true, AutoHeartBeat: true, HeartBeatInterval: time.Second}

	client.OnClose = func(c *client2.Client) {
		console.Info("udp client close")
	}

	client.OnOpen = func(c *client2.Client) {
		console.Info("udp client open")
	}

	client.OnError = func(err error) {
		console.Info("udp client", err)
	}

	// create router
	var router = &client2.Router{IgnoreCase: true}

	Router(router)

	client.OnSuccess = func() {
		console.Info("udp client start success")
		fn()
	}

	app.Node.Client = client

	go client.SetRouter(router).Connect()
}
