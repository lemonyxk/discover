/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-04 19:07
**/

package tcp

import (
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/v2"
	"github.com/lemonyxk/kitty/v2/socket/websocket/server"
)

func Start(host string, fn func()) {

	var tcpServer = server.Server{Name: host, Addr: host, HeartBeatTimeout: 3 * time.Second}

	tcpServer.OnClose = func(conn server.Conn) {

		app.Node.Lock()

		defer app.Node.Unlock()

		console.Info("tcp server", conn.FD(), "close")

		var data = app.Node.Register.Get(conn.FD())
		if data == nil {
			return
		}

		app.Node.Register.Delete(conn.FD())

		for i := 0; i < len(data.ServerList); i++ {
			app.Node.Alive.DeleteConn(data.ServerList[i], conn)
		}

		for i := 0; i < len(data.KeyList); i++ {
			app.Node.Key.Delete(data.KeyList[i], conn)
		}

		if data.ServerInfo != nil {
			app.Node.Alive.DeleteData(data.ServerInfo.ServerName, data.ServerInfo.Addr)
			var list = app.Node.Alive.GetData(data.ServerInfo.ServerName)
			var connections = app.Node.Alive.GetConn(data.ServerInfo.ServerName)
			for i := 0; i < len(connections); i++ {
				var err = connections[i].ProtoBufEmit("/Alive", &message.ServerInfoList{List: list})
				if err != nil {
					console.Error(err)
				}
			}
		}

	}

	tcpServer.OnError = func(err error) {
		console.Error("tcp server", err)
	}

	tcpServer.OnOpen = func(conn server.Conn) {
		console.Info("tcp server", conn.FD(), "open")
	}

	var router = kitty.NewWebSocketServerRouter()

	Router(router)

	tcpServer.OnSuccess = func() {
		console.Info("tcp server start success", tcpServer.LocalAddr())
		fn()
	}

	app.Node.Server = &tcpServer

	go tcpServer.SetRouter(router).Start()

}
