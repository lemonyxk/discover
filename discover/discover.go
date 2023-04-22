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
	"fmt"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/store"
	"github.com/lemonyxk/kitty/errors"
	"github.com/lemonyxk/kitty/socket"
	client2 "github.com/lemonyxk/kitty/socket/websocket/client"
)

var hasAlive int32 = 0
var hasKey int32 = 0

type Client struct {
	serverList []*message.Address
	master     *message.Server
	register   *client2.Client
	listen     *client2.Client

	registerFn func()
	aliveFn    func()
	listenFn   func()
}

type Alive struct {
	client *Client
	list   []string
}

func (dis *Client) Register(serverName, addr string) {

	if serverName == "" || addr == "" {
		panic("server name or addr is empty")
	}

	dis.registerFn = func() {
		var bts, err = jsoniter.Marshal(&message.ServerInfo{
			Name: serverName,
			Addr: addr,
		})
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			dis.registerFn()
			return
		}
		err = dis.register.Sender().Emit("/Register", bts)
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			dis.registerFn()
			return
		}
	}

	dis.registerFn()
}

func (dis *Client) Alive(serverList ...string) *Alive {

	if atomic.AddInt32(&hasAlive, 1) > 1 {
		panic("repeat monitoring")
	}

	return &Alive{
		client: dis,
		list:   serverList,
	}
}

func (w *Alive) Watch(fn func(data []*message.ServerInfo)) {

	if len(w.list) == 0 {
		return
	}

	w.client.aliveFn = func() {
		w.client.register.GetRouter().Remove("/Alive")
		w.client.register.GetRouter().Route("/Alive").Handler(func(stream *socket.Stream[client2.Conn]) error {
			if stream.Code() != 200 {
				return errors.New(fmt.Sprintf("alive error:%d %s", stream.Code(), stream.Data))
			}
			var res []*message.ServerInfo
			var err = jsoniter.Unmarshal(stream.Data, &res)
			if err != nil {
				return errors.New(err)
			}
			fn(res)
			return nil
		})

		var bts, err = jsoniter.Marshal(w.list)
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			w.client.aliveFn()
			return
		}

		err = w.client.register.Sender().Emit("/Alive", bts)
		if err != nil {
			time.Sleep(time.Second)
			w.client.aliveFn()
			return
		}
	}

	w.client.aliveFn()
}

type KeyList struct {
	dis  *Client
	list []string
}

func (dis *Client) Key(keyList ...string) *KeyList {

	if atomic.AddInt32(&hasKey, 1) > 1 {
		panic("repeat monitoring")
	}

	return &KeyList{
		dis:  dis,
		list: keyList,
	}
}

func (k *KeyList) Watch(fn func(op *store.Message)) {

	if len(k.list) == 0 {
		return
	}

	k.dis.listenFn = func() {
		k.dis.listen.GetRouter().Remove("/Key")
		k.dis.listen.GetRouter().Route("/Key").Handler(func(stream *socket.Stream[client2.Conn]) error {
			if stream.Code() != 200 {
				return errors.New(stream.Code())
			}
			msg, err := store.Parse(stream.Data)
			if err != nil {
				return errors.New(err)
			}
			fn(msg)
			return nil
		})

		var bts, err = jsoniter.Marshal(k.list)
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			k.dis.listenFn()
			return
		}

		err = k.dis.listen.Sender().Emit("/Key", bts)
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			k.dis.listenFn()
			return
		}
	}

	k.dis.listenFn()
}

func (dis *Client) refreshMaster() {
	var register = dis.getMasterServer()
	dis.register.Addr = "ws://" + register.Tcp
	console.Info("new register addr:", register.Addr)
}

func (dis *Client) refreshCluster() {
	var listen = dis.randomAddr()
	dis.listen.Addr = "ws://" + listen.Tcp
	console.Info("new listen addr:", listen.Addr)
}
