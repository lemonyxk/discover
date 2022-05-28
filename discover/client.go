/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-27 15:39
**/

package discover

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/v2"
	client2 "github.com/lemonyxk/kitty/v2/socket/websocket/client"
)

func New(serverList []string) *discover {

	if len(serverList) == 0 {
		panic("server list is empty")
	}

	var dis = &discover{}

	for i := 0; i < len(serverList); i++ {
		dis.serverList = append(dis.serverList, &message.WhoIsMaster{Addr: app.ParseAddr(serverList[i])})
	}

	dis.getMasterServer()

	var wait sync.WaitGroup

	wait.Add(2)

	initRegister(dis, &wait)

	initLister(dis, &wait)

	wait.Wait()

	return dis
}

func initRegister(dis *discover, wait *sync.WaitGroup) {

	var isStart int32 = 0

	var client = &client2.Client{
		Addr:              "ws://" + dis.master.Tcp,
		HeartBeatTimeout:  3 * time.Second,
		HeartBeatInterval: 1 * time.Second,
		ReconnectInterval: 1 * time.Second,
	}

	dis.register = client

	dis.register.OnOpen = func(conn client2.Conn) {
		console.Info("register client open at:", dis.register.Addr)
	}

	dis.register.OnClose = func(conn client2.Conn) {
		console.Info("register client close")
	}

	dis.register.OnError = func(err error) {
		console.Info("register client error:", err)
	}

	dis.register.OnReconnecting = func() {
		dis.refreshMaster()
	}

	dis.register.OnSuccess = func() {
		if atomic.AddInt32(&isStart, 1) == 1 {
			wait.Add(-1)
		}
		if dis.registerFn != nil {
			go dis.registerFn()
		}
		if dis.aliveFn != nil {
			go dis.aliveFn()
		}
	}

	var r = kitty.NewWebSocketClientRouter()

	dis.register.SetRouter(r)

	go dis.register.Connect()
}

func initLister(dis *discover, wait *sync.WaitGroup) {

	var isStart int32 = 0

	var client = &client2.Client{
		Addr:              "ws://" + dis.randomAddr().Tcp,
		HeartBeatTimeout:  3 * time.Second,
		HeartBeatInterval: 1 * time.Second,
		ReconnectInterval: 1 * time.Second,
	}

	dis.listen = client

	dis.listen.OnOpen = func(conn client2.Conn) {
		console.Info("listen client open at:", dis.listen.Addr)
	}

	dis.listen.OnClose = func(conn client2.Conn) {
		console.Info("listen client close")
	}

	dis.listen.OnError = func(err error) {
		console.Info("listen client error:", err)
	}

	dis.listen.OnReconnecting = func() {
		dis.refreshCluster()
	}

	dis.listen.OnSuccess = func() {
		if atomic.AddInt32(&isStart, 1) == 1 {
			wait.Add(-1)
		}
		if dis.listenFn != nil {
			go dis.listenFn()
		}
	}

	var r = kitty.NewWebSocketClientRouter()

	dis.listen.SetRouter(r)

	go dis.listen.Connect()
}
