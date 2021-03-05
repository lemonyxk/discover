/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-25 18:36
**/

package server

import (
	"time"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/kitty/socket/udp"
	"github.com/lemoyxk/kitty/socket/udp/server"

	"discover/app"
)

func Start(host string, fn func()) {
	// create server
	var udpServer = &server.Server{Name: host, Host: host, HeartBeatTimeout: 3 * time.Second}

	// event
	udpServer.OnClose = func(conn *server.Conn) {
		console.Info("udp server", conn.FD, "close")
	}

	udpServer.OnError = func(err error) {
		console.Error("udp server", err)
	}

	udpServer.OnOpen = func(conn *server.Conn) {
		console.Info("udp server", conn.FD, "open")
		for {
			// if start too early
			// client is not ready
			if app.Node.Client == nil || app.Node.Client.Conn == nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			// tell others you are ready
			_ = app.Node.Client.Push(udp.OpenMessage)
			break
		}
	}

	// create router
	var router = server.Router{IgnoreCase: true}

	Router(&router)

	udpServer.OnSuccess = func() {
		console.Info("udp server start success", udpServer.LocalAddr())
		fn()
	}

	go udpServer.SetRouter(&router).Start()

}
