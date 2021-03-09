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

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/websocket/server"

	"discover/app"
	"discover/message"
)

func Start(host string, fn func()) {

	var tcpServer = server.Server{Name: host, Host: host, HeartBeatTimeout: 3 * time.Second}

	tcpServer.OnClose = func(conn *server.Conn) {

		app.Node.Lock()

		defer app.Node.Unlock()

		console.Info("tcp server", conn.FD, "close")

		var data = app.Node.Register.Get(conn.FD)
		if data == nil {
			return
		}

		app.Node.Register.Delete(conn.FD)

		for i := 0; i < len(data.ServerList); i++ {
			app.Node.Alive.DeleteConn(data.ServerList[i], conn)
		}

		for i := 0; i < len(data.KeyList); i++ {
			app.Node.Listen.Delete(data.KeyList[i], conn)
		}

		if data.ServerInfo != nil {
			app.Node.Alive.DeleteData(data.ServerInfo.ServerName, data.ServerInfo.Addr)
			var list = app.Node.Alive.GetData(data.ServerInfo.ServerName)
			var connections = app.Node.Alive.GetConn(data.ServerInfo.ServerName)
			for i := 0; i < len(connections); i++ {
				_ = connections[i].ProtoBufEmit(socket.ProtoBufPack{
					Event: "/OnRegister",
					Data:  &message.ServerInfoList{List: list},
				})
			}
		}

	}

	tcpServer.OnError = func(err error) {
		console.Error("tcp server", err)
	}

	tcpServer.OnOpen = func(conn *server.Conn) {
		console.Info("tcp server", conn.FD, "open")
	}

	var router = server.Router{IgnoreCase: true}

	Router(&router)

	tcpServer.OnSuccess = func() {
		console.Info("tcp server start success", tcpServer.LocalAddr())
		fn()
	}

	app.Node.Server = &tcpServer

	go tcpServer.SetRouter(&router).Start()

}
