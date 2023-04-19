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
	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/structs"
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

	var conn = stream.Conn()

	var data message.ServerInfo

	var err = jsoniter.Unmarshal(stream.Data, &data)
	if err != nil {
		return api.Failed(conn, stream.Event, err.Error())
	}

	if data.Name == "" || data.Addr == "" {
		return api.Failed(conn, stream.Event, "server name or address is empty")
	}

	var register = app.Node.Register.Get(conn.FD())
	if register == nil {
		register = &structs.Register{}
	}
	register.ServerInfo = &data

	app.Node.Register.Set(conn.FD(), register)

	// add to watch queue
	app.Node.Alive.AddData(data.Name, data.Addr)
	var list = app.Node.Alive.GetData(data.Name)

	var connections = app.Node.Alive.GetConn(data.Name)
	for i := 0; i < len(connections); i++ {
		var err = api.Success(connections[i], "/Alive", list)
		if err != nil {
			console.Error(err)
		}
	}

	return nil
}

func (api *action) Alive(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var conn = stream.Conn()

	var list []string

	var err = jsoniter.Unmarshal(stream.Data, &list)
	if err != nil {
		return api.Failed(conn, stream.Event, err.Error())
	}

	if len(list) == 0 {
		return api.Failed(conn, stream.Event, "server list is empty")
	}

	var register = app.Node.Register.Get(conn.FD())
	if register == nil {
		register = &structs.Register{}
	}

	register.ServerList = list

	app.Node.Register.Set(conn.FD(), register)

	// add to notify queue
	for i := 0; i < len(list); i++ {
		app.Node.Alive.AddConn(list[i], conn)
	}

	// notify what you are watching
	for i := 0; i < len(list); i++ {
		var list = app.Node.Alive.GetData(list[i])
		if len(list) == 0 {
			continue
		}
		var err = api.Success(conn, "/Alive", list)
		if err != nil {
			console.Error(err)
		}
	}

	return nil
}

func (api *action) Key(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()
	defer app.Node.Unlock()

	var conn = stream.Conn()

	var list []string

	var err = jsoniter.Unmarshal(stream.Data, &list)
	if err != nil {
		return api.Failed(conn, stream.Event, err.Error())
	}

	if len(list) == 0 {
		return api.Failed(conn, stream.Event, "key list is empty")
	}

	var register = app.Node.Register.Get(conn.FD())
	if register == nil {
		register = &structs.Register{}
	}

	register.KeyList = list

	app.Node.Register.Set(conn.FD(), register)

	// add to watch queue
	for i := 0; i < len(list); i++ {

		var key = list[i]
		app.Node.Key.Add(key, conn)

		var value, err = app.Node.Store.Get(key)
		if err != nil {
			continue
		}

		if value == "" {
			continue
		}

		err = api.Success(conn, stream.Event, &message.Op{Key: key, Value: value, Op: "set"})
		if err != nil {
			console.Error(err)
		}
	}

	return nil
}
