/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-27 15:39
**/

package discover

import (
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/message"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/kitty"
	"github.com/lemonyxk/kitty/socket"
	client2 "github.com/lemonyxk/kitty/socket/websocket/client"
)

type Discover struct {
	config     *Config
	serverList []string
}

func (d *Discover) Config(config *Config) *Discover {
	d.config = config
	return d
}

func (d *Discover) Connect() *Client {
	var client = &Client{
		registerClose: make(chan struct{}, 1),
		listenClose:   make(chan struct{}, 1),
		config:        d.config,
	}

	for i := 0; i < len(d.serverList); i++ {
		client.serverList = append(client.serverList, &message.Address{
			Server: app.ParseAddr(d.serverList[i]),
		})
	}

	client.getMasterServer()

	var wait sync.WaitGroup

	wait.Add(2)

	initRegister(client, &wait)

	initLister(client, &wait)

	wait.Wait()

	return client
}

func New(serverList ...string) *Discover {
	if len(serverList) == 0 {
		panic("server list is empty")
	}

	return &Discover{
		serverList: serverList,
	}
}

func initRegister(dis *Client, wait *sync.WaitGroup) {

	var isStart int32 = 0

	var client = &client2.Client{
		Addr:              "ws://" + dis.master.Tcp,
		HeartBeatTimeout:  6 * time.Second,
		HeartBeatInterval: 1 * time.Second,
		ReconnectInterval: 1 * time.Second,
	}

	dis.register = client

	dis.register.OnOpen = func(conn client2.Conn) {
		console.Errorf("register client open at: %s\n", dis.register.Addr)
	}

	dis.register.OnClose = func(conn client2.Conn) {
		console.Infof("register client close at: %s\n", dis.register.Addr)
	}

	dis.register.OnError = func(stream *socket.Stream[client2.Conn], err error) {
		console.Infof("register client error: %+v\n", err)
	}

	dis.register.OnException = func(err error) {
		console.Infof("register client exception: %+v\n", err)
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

func initLister(dis *Client, wait *sync.WaitGroup) {

	var isStart int32 = 0

	var client = &client2.Client{
		Addr:              "ws://" + dis.randomAddr().Tcp,
		HeartBeatTimeout:  3 * time.Second,
		HeartBeatInterval: 1 * time.Second,
		ReconnectInterval: 1 * time.Second,
	}

	dis.listen = client

	dis.listen.OnOpen = func(conn client2.Conn) {
		console.Infof("listen client open at: %s\n", dis.listen.Addr)
	}

	dis.listen.OnClose = func(conn client2.Conn) {
		console.Infof("listen client close at: %s\n", dis.listen.Addr)
	}

	dis.listen.OnError = func(stream *socket.Stream[client2.Conn], err error) {
		console.Errorf("listen client error: %+v\n", err)
	}

	dis.listen.OnException = func(err error) {
		console.Errorf("listen client exception: %+v\n", err)
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
