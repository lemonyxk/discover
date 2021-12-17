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
	"github.com/lemoyxk/discover/app"
	"github.com/lemoyxk/kitty"
	client2 "github.com/lemoyxk/kitty/socket/udp/client"
)

func Start(host string, fn func()) {

	var client = &client2.Client{
		Addr:              host,
		HeartBeatTimeout:  3 * time.Second,
		HeartBeatInterval: 1 * time.Second,
		ReconnectInterval: 1 * time.Second,
	}

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
	var router = kitty.NewUdpClientRouter()

	Router(router)

	client.OnSuccess = func() {
		console.Info("udp client start success")
		fn()
	}

	app.Node.Client = client

	go client.SetRouter(router).Connect()
}
