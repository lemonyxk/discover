/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-27 16:35
**/

package discover

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/v2/socket"
	client2 "github.com/lemonyxk/kitty/v2/socket/websocket/client"
	"google.golang.org/protobuf/proto"
)

var hasAlive int32 = 0
var hasKey int32 = 0

type discover struct {
	serverList []*message.WhoIsMaster
	master     *message.Address
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
		var err = dis.register.ProtoBufEmit("/Register", &message.ServerInfo{
			ServerName: serverName,
			Addr:       addr,
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
			var data message.ServerInfoList
			var err = proto.Unmarshal(stream.Data, &data)
			if err != nil {
				return err
			}
			fn(data.List)
			return nil
		})

		var err = w.dis.register.ProtoBufEmit("/Alive", &message.ServerList{List: w.serverList})
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

func (k *key) Watch(fn func(key, value string)) {

	if len(k.keyList) == 0 {
		return
	}

	k.dis.listenFn = func() {
		k.dis.listen.GetRouter().Remove("/Key")
		k.dis.listen.GetRouter().Route("/Key").Handler(func(stream *socket.Stream[client2.Conn]) error {
			var data = string(stream.Data)
			var index = strings.Index(data, "\n")
			fn(data[:index], data[index+1:])
			return nil
		})

		var err = k.dis.listen.ProtoBufEmit("/Key", &message.KeyList{List: k.keyList})
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
