/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-04 19:07
**/

package tcp

import (
	json "github.com/bytedance/sonic"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty"
	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

func Start(host string, fn func()) {

	var tcpServer = server.Server[any]{Name: host, Addr: host, HeartBeatTimeout: 6 * time.Second}

	tcpServer.OnClose = func(conn server.Conn) {

		app.Node.Lock()

		defer app.Node.Unlock()

		console.Info("tcp server", conn.FD(), "close")

		var data = app.Node.Register.Get(conn.FD())
		if data == nil {
			return
		}

		app.Node.Register.Delete(conn.FD())

		console.Info("tcp server", conn.FD(), "unregister", data.ServerInfo)

		for i := 0; i < len(data.ServerList); i++ {
			app.Node.Alive.DeleteConn(data.ServerList[i], conn.FD())
		}

		for i := 0; i < len(data.KeyList); i++ {
			app.Node.Key.Delete(data.KeyList[i], conn.FD())
		}

		if data.ServerInfo != nil {
			app.Node.Alive.DeleteData(data.ServerInfo.Name, data.ServerInfo.Addr)
			var list = app.Node.Alive.GetData(data.ServerInfo.Name)
			var connections = app.Node.Alive.GetConn(data.ServerInfo.Name)
			for i := 0; i < len(connections); i++ {
				connections[i].SetCode(200)
				var bts, err = json.Marshal(message.AliveResponse{Name: data.ServerInfo.Name, ServerInfoList: list})
				if err != nil {
					console.Error(err)
					continue
				}
				err = connections[i].Emit("/Alive", bts)
				if err != nil {
					console.Error(err)
					continue
				}
			}
		}
	}

	tcpServer.OnError = func(stream *socket.Stream[server.Conn], err error) {
		console.Error("tcp server", err)
	}

	tcpServer.OnOpen = func(conn server.Conn) {
		console.Info("tcp server", conn.FD(), "open")
	}

	var router = kitty.NewWebSocketServerRouter[any]()

	Router(router)

	tcpServer.OnSuccess = func() {
		console.Info("tcp server start at", host, "success")
		fn()
	}

	app.Node.Server = &tcpServer

	go tcpServer.SetRouter(router).Start()

}
