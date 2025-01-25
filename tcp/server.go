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
	json "github.com/lemonyxk/kitty/json"
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

		console.Info.Logf("tcp server %d close", conn.FD())

		var data = app.Node.Register.Get(conn.FD())
		if data == nil {
			return
		}

		app.Node.Register.Delete(conn.FD())

		console.Info.Logf("tcp server %d unregister %v", conn.FD(), data.ServerInfo)

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
					console.Error.Logf("%+v", err)
					continue
				}
				err = connections[i].Emit("/Alive", bts)
				if err != nil {
					console.Error.Logf("%+v", err)
					continue
				}
			}
		}
	}

	tcpServer.OnError = func(stream *socket.Stream[server.Conn], err error) {
		console.Error.Logf("tcp server error: %s", err)
	}

	tcpServer.OnOpen = func(conn server.Conn) {
		console.Info.Logf("tcp server %d open", conn.FD())
	}

	var router = kitty.NewWebSocketServerRouter[any]()

	Router(router)

	tcpServer.OnSuccess = func() {
		console.Info.Logf("tcp server start at: %s success", host)
		fn()
	}

	app.Node.Server = &tcpServer

	go tcpServer.SetRouter(router).Start()

}
