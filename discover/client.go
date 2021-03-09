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

	"github.com/lemoyxk/console"
	client2 "github.com/lemoyxk/kitty/socket/websocket/client"

	"discover/app"
	"discover/message"
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
		Scheme:            "ws",
		Host:              dis.master.Tcp,
		Reconnect:         true,
		AutoHeartBeat:     true,
		HeartBeatTimeout:  3 * time.Second,
		HeartBeatInterval: time.Second,
	}

	dis.register = client

	dis.register.OnOpen = func(c *client2.Client) {
		console.Info("register client open at:", dis.register.Host)
	}

	dis.register.OnClose = func(c *client2.Client) {
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

	var r = &client2.Router{IgnoreCase: true}

	dis.register.SetRouter(r)

	go dis.register.Connect()
}

func initLister(dis *discover, wait *sync.WaitGroup) {

	var isStart int32 = 0

	var client = &client2.Client{
		Scheme:            "ws",
		Host:              dis.randomAddr().Tcp,
		Reconnect:         true,
		AutoHeartBeat:     true,
		HeartBeatTimeout:  3 * time.Second,
		HeartBeatInterval: time.Second,
	}

	dis.listen = client

	dis.listen.OnOpen = func(c *client2.Client) {
		console.Info("listen client open at:", dis.listen.Host)
	}

	dis.listen.OnClose = func(c *client2.Client) {
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

	var r = &client2.Router{IgnoreCase: true}

	dis.listen.SetRouter(r)

	go dis.listen.Connect()
}
