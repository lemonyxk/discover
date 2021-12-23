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
	"errors"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/discover/message"
	"github.com/lemoyxk/kitty/socket"
	client2 "github.com/lemoyxk/kitty/socket/websocket/client"
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

		var stream, err = dis.register.Async().ProtoBufEmit(socket.ProtoBufPack{
			Event: "/Register",
			Data: &message.ServerInfo{
				ServerName: serverName,
				Addr:       addr,
			},
		})
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			dis.registerFn()
			return
		}

		if string(stream.Data) != "OK" {
			console.Info(errors.New(string(stream.Data)))
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
		w.dis.register.GetRouter().Remove("/OnRegister")
		w.dis.register.GetRouter().Route("/OnRegister").Handler(func(client *client2.Client, stream *socket.Stream) error {
			var data message.ServerInfoList
			var err = proto.Unmarshal(stream.Data, &data)
			if err != nil {
				return err
			}
			fn(data.List)
			return nil
		})

		var err = w.dis.register.ProtoBufEmit(socket.ProtoBufPack{
			Event: "/OnRegister",
			Data:  &message.ServerList{List: w.serverList},
		})
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

		k.dis.listen.GetRouter().Remove("/OnListen")
		k.dis.listen.GetRouter().Route("/OnListen").Handler(func(client *client2.Client, stream *socket.Stream) error {
			var data = string(stream.Data)
			var index = strings.Index(data, "\n")
			fn(data[:index], data[index+1:])
			return nil
		})

		var stream, err = k.dis.listen.Async().ProtoBufEmit(socket.ProtoBufPack{
			Event: "/Listen",
			Data:  &message.KeyList{List: k.keyList},
		})
		if err != nil {
			console.Info(err)
			time.Sleep(time.Second)
			k.dis.listenFn()
			return
		}

		if string(stream.Data) != "OK" {
			console.Info(errors.New(string(stream.Data)))
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
