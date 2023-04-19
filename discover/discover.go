/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-27 16:35
**/

package discover

import (
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/errors"
	"github.com/lemonyxk/kitty/socket"
	client2 "github.com/lemonyxk/kitty/socket/websocket/client"
)

var hasAlive int32 = 0
var hasKey int32 = 0

type discover struct {
	serverList []*message.Address
	master     *message.Server
	register   *client2.Client
	listen     *client2.Client

	registerFn func()
	aliveFn    func()
	listenFn   func()
}

type alive struct {
	dis        *discover
	serverList []string
}

func (dis *discover) Register(serverName, addr string) {

	if serverName == "" || addr == "" {
		panic("server name or addr is empty")
	}

	dis.registerFn = func() {
		var err = dis.register.JsonEmit("/Register", &message.ServerInfo{
			Name: serverName,
			Addr: addr,
		})
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			dis.registerFn()
			return
		}
	}

	dis.registerFn()
}

func (dis *discover) Alive(serverList ...string) *alive {

	if atomic.AddInt32(&hasAlive, 1) > 1 {
		panic("repeat monitoring")
	}

	return &alive{
		dis:        dis,
		serverList: serverList,
	}
}

func (w *alive) Watch(fn func(data []*message.ServerInfo)) {

	if len(w.serverList) == 0 {
		return
	}

	w.dis.aliveFn = func() {
		w.dis.register.GetRouter().Remove("/Alive")
		w.dis.register.GetRouter().Route("/Alive").Handler(func(stream *socket.Stream[client2.Conn]) error {
			var res message.ServerInfoResponse
			var err = jsoniter.Unmarshal(stream.Data, &res)
			if err != nil {
				return errors.New(err)
			}
			if res.Code != 200 {
				return errors.New(res.Code)
			}
			fn(res.Msg)
			return nil
		})

		var err = w.dis.register.JsonEmit("/Alive", w.serverList)
		if err != nil {
			time.Sleep(time.Second)
			w.dis.aliveFn()
			return
		}
	}

	w.dis.aliveFn()
}

type key struct {
	dis     *discover
	keyList []string
}

func (dis *discover) Key(keyList ...string) *key {

	if atomic.AddInt32(&hasKey, 1) > 1 {
		panic("repeat monitoring")
	}

	return &key{
		dis:     dis,
		keyList: keyList,
	}
}

func (k *key) Watch(fn func(op message.Op)) {

	if len(k.keyList) == 0 {
		return
	}

	k.dis.listenFn = func() {
		k.dis.listen.GetRouter().Remove("/Key")
		k.dis.listen.GetRouter().Route("/Key").Handler(func(stream *socket.Stream[client2.Conn]) error {
			var res message.OpResponse
			var err = jsoniter.Unmarshal(stream.Data, &res)
			if err != nil {
				return errors.New(err)
			}
			if res.Code != 200 {
				return errors.New(res.Code)
			}
			fn(res.Msg)
			return nil
		})

		var err = k.dis.listen.JsonEmit("/Key", k.keyList)
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			k.dis.listenFn()
			return
		}
	}

	k.dis.listenFn()
}

func (dis *discover) refreshMaster() {
	var register = dis.getMasterServer()
	dis.register.Addr = "ws://" + register.Tcp
	console.Info("new register addr:", register.Addr)
}

func (dis *discover) refreshCluster() {
	var listen = dis.randomAddr()
	dis.listen.Addr = "ws://" + listen.Tcp
	console.Info("new listen addr:", listen.Addr)
}
