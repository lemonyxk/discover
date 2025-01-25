/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-04 19:15
**/

package tcp

import (
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/store"
	"github.com/lemonyxk/discover/structs"
	json "github.com/lemonyxk/kitty/json"
	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

var Action = &action{}

type action struct {
	Controller
}

func (api *action) Register(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()
	defer app.Node.Unlock()

	var data message.ServerInfo

	var err = json.Unmarshal(stream.Data(), &data)
	if err != nil {
		return api.Failed(stream, stream.Event(), err.Error())
	}

	if data.Name == "" || data.Addr == "" {
		return api.Failed(stream, stream.Event(), "server name or address is empty")
	}

	// add to watch queue
	if !app.Node.Alive.AddData(data) {
		return api.Failed(stream, stream.Event(), "server name already exists")
	}

	var register = app.Node.Register.Get(stream.Conn().FD())
	if register == nil {
		register = &structs.Register{}
	}
	register.ServerInfo = &data

	app.Node.Register.Set(stream.Conn().FD(), register)

	var list = app.Node.Alive.GetData(data.Name)

	var connections = app.Node.Alive.GetConn(data.Name)
	for i := 0; i < len(connections); i++ {
		connections[i].SetCode(200)
		var bts, err = json.Marshal(message.AliveResponse{Name: data.Name, ServerInfoList: list})
		if err != nil {
			console.Error.Logf("%+v", err)
			continue
		}
		err = connections[i].Emit("/Alive", bts)
		if err != nil {
			console.Error.Logf("%+v", err)
		}
	}

	console.Info.Logf("tcp server %d register %s %s", stream.Conn().FD(), data.Name, data.Addr)

	stream.SetCode(200)

	return stream.Emit("/Register", nil)
}

func (api *action) Update(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()
	defer app.Node.Unlock()

	var data message.ServerInfo

	var err = json.Unmarshal(stream.Data(), &data)
	if err != nil {
		return api.Failed(stream, stream.Event(), err.Error())
	}

	if data.Name == "" || data.Addr == "" {
		return api.Failed(stream, stream.Event(), "server name or address is empty")
	}

	// add to watch queue
	if !app.Node.Alive.SetData(data) {
		return api.Failed(stream, stream.Event(), "server name not exists")
	}

	var list = app.Node.Alive.GetData(data.Name)

	var connections = app.Node.Alive.GetConn(data.Name)
	for i := 0; i < len(connections); i++ {
		connections[i].SetCode(200)
		var bts, err = json.Marshal(message.AliveResponse{Name: data.Name, ServerInfoList: list})
		if err != nil {
			console.Error.Logf("%+v", err)
			continue
		}
		err = connections[i].Emit("/Alive", bts)
		if err != nil {
			console.Error.Logf("%+v", err)
		}
	}

	console.Info.Logf("tcp server %d update %s %s", stream.Conn().FD(), data.Name, data.Addr)

	stream.SetCode(200)

	return stream.Emit("/Update", nil)
}

func (api *action) Alive(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var list []string

	var err = json.Unmarshal(stream.Data(), &list)
	if err != nil {
		return api.Failed(stream, stream.Event(), err.Error())
	}

	if len(list) == 0 {
		return api.Failed(stream, stream.Event(), "server list is empty")
	}

	var register = app.Node.Register.Get(stream.Conn().FD())
	if register == nil {
		register = &structs.Register{}
	}

	register.ServerList = list

	app.Node.Register.Set(stream.Conn().FD(), register)

	// add to notify queue
	for i := 0; i < len(list); i++ {
		app.Node.Alive.AddConn(list[i], stream)
	}

	// notify what you are watching
	for i := 0; i < len(list); i++ {
		var info = app.Node.Alive.GetData(list[i])
		if len(info) == 0 {
			continue
		}

		stream.SetCode(200)
		var bts, err = json.Marshal(message.AliveResponse{Name: list[i], ServerInfoList: info})
		if err != nil {
			console.Error.Logf("%+v", err)
			continue
		}
		err = api.Success(stream, "/Alive", bts)
		if err != nil {
			console.Error.Logf("%+v", err)
		}
	}

	return nil
}

func (api *action) Key(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()
	defer app.Node.Unlock()

	var list []string

	var err = json.Unmarshal(stream.Data(), &list)
	if err != nil {
		return api.Failed(stream, stream.Event(), err.Error())
	}

	if len(list) == 0 {
		return api.Failed(stream, stream.Event(), "key list is empty")
	}

	var register = app.Node.Register.Get(stream.Conn().FD())
	if register == nil {
		register = &structs.Register{}
	}

	register.KeyList = list

	app.Node.Register.Set(stream.Conn().FD(), register)

	// add to watch queue
	for i := 0; i < len(list); i++ {

		var key = list[i]
		app.Node.Key.Add(key, stream)

		var value, err = app.Node.Store.Get(key)
		if err != nil {
			continue
		}

		if len(value) == 0 {
			continue
		}

		msg := store.Build(&store.Message{Op: store.Set, Key: key, Value: value})
		err = api.Success(stream, stream.Event(), msg)
		if err != nil {
			console.Error.Logf("%+v", err)
		}
	}

	return nil
}
