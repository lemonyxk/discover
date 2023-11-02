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

type Config struct {
	AutoUpdateInterval time.Duration
}

var hasAlive int32 = 0
var hasKey int32 = 0

type Client struct {
	config *Config

	serverList []*message.Address
	master     *message.Server
	register   *client2.Client[any]
	listen     *client2.Client[any]

	registerFn func()
	aliveFn    func()
	listenFn   func()

	registerClose chan struct{}
	listenClose   chan struct{}
}

type Alive struct {
	client  *Client
	list    []string
	onClose func()
}

func (dis *Client) Register(fn func() message.ServerInfo) {

	var info = fn()

	if info.Name == "" || info.Addr == "" {
		panic("server name or addr is empty")
	}

	dis.registerFn = func() {
		dis.register.GetRouter().Remove("/Register")
		dis.register.GetRouter().Route("/Register").Handler(func(stream *socket.Stream[client2.Conn]) error {
			if stream.Code() != 200 {
				if string(stream.Data()) == "NOT MASTER" {
					if dis.register != nil {
						_ = dis.register.Close()
					}
					if dis.listen != nil {
						_ = dis.listen.Close()
					}
				}
				return errors.New(fmt.Sprintf("register error:%d %s", stream.Code(), stream.Data()))
			}
			go func() {
				if dis.config == nil || dis.config.AutoUpdateInterval == 0 {
					return
				}

				var ticker = time.NewTicker(dis.config.AutoUpdateInterval)

				dis.register.GetRouter().Remove("/Update")
				dis.register.GetRouter().Route("/Update").Handler(func(stream *socket.Stream[client2.Conn]) error {
					if stream.Code() != 200 {
						if string(stream.Data()) == "NOT MASTER" {
							_ = dis.register.Close()
						}
						return errors.New(fmt.Sprintf("update error:%d %s", stream.Code(), stream.Data()))
					}
					return nil
				})

				for {
					select {
					case <-dis.registerClose:
						ticker.Stop()
						return
					case <-ticker.C:
						var info = fn()
						var bts, err = jsoniter.Marshal(info)
						if err != nil {
							console.Info(err)
							continue
						}
						err = dis.register.Sender().Emit("/Update", bts)
						if err != nil {
							console.Info(err)
							continue
						}
					}
				}

			}()

			return nil
		})
		var bts, err = jsoniter.Marshal(info)
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

func (w *Alive) Watch(fn func(name string, serverInfo []*message.ServerInfo)) {

	if len(w.list) == 0 {
		return
	}

	w.client.aliveFn = func() {
		w.client.register.GetRouter().Remove("/Alive")
		w.client.register.GetRouter().Route("/Alive").Handler(func(stream *socket.Stream[client2.Conn]) error {
			if stream.Code() != 200 {
				if string(stream.Data()) == "NOT MASTER" {
					if w.client.register != nil {
						_ = w.client.register.Close()
					}
					if w.client.listen != nil {
						_ = w.client.listen.Close()
					}
				}
				return errors.New(fmt.Sprintf("alive error:%d %s", stream.Code(), stream.Data()))
			}
			var res message.AliveResponse
			var err = jsoniter.Unmarshal(stream.Data(), &res)
			if err != nil {
				return errors.New(err)
			}
			fn(res.Name, res.ServerInfoList)
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

func (w *Alive) OnClose(fn func()) {
	w.client.register.OnClose = func(conn client2.Conn) {
		w.client.registerClose <- struct{}{}
		fn()
	}
}

type KeyList struct {
	dis     *Client
	list    []string
	onClose func()
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
				return errors.New(fmt.Sprintf("alive error:%d %s", stream.Code(), stream.Data()))
			}
			msg, err := store.Parse(stream.Data())
			if err != nil {
				return errors.New(err)
			}
			fn(msg)

			go func() {
				<-k.dis.listenClose // wait listen close and do something
			}()

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

func (k *KeyList) OnClose(fn func()) {
	k.dis.listen.OnClose = func(conn client2.Conn) {
		k.dis.listenClose <- struct{}{}
		fn()
	}
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
