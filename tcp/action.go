/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-04 19:15
**/

package tcp

import (
	"errors"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/discover/structs"
	"github.com/lemonyxk/kitty/v2/socket"
	"github.com/lemonyxk/kitty/v2/socket/websocket/server"
	"google.golang.org/protobuf/proto"
)

func Register(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()
	defer app.Node.Unlock()

	var conn = stream.Conn

	var data message.ServerInfo

	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	if data.ServerName == "" || data.Addr == "" {
		return errors.New("server name or addr is empty")
	}

	var register = app.Node.Register.Get(conn.FD())
	if register == nil {
		register = &structs.Register{}
	}
	register.ServerInfo = &data

	app.Node.Register.Set(conn.FD(), register)

	// add to watch queue
	app.Node.Alive.AddData(data.ServerName, data.Addr)
	var list = app.Node.Alive.GetData(data.ServerName)

	var connections = app.Node.Alive.GetConn(data.ServerName)
	for i := 0; i < len(connections); i++ {
		var err = connections[i].ProtoBufEmit("/Alive", &message.ServerInfoList{List: list})
		if err != nil {
			console.Error(err)
		}
	}

	return nil
}

func Alive(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var conn = stream.Conn

	var data message.ServerList

	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	if len(data.List) == 0 {
		return errors.New("server list is empty")
	}

	var register = app.Node.Register.Get(conn.FD())
	if register == nil {
		register = &structs.Register{}
	}
	register.ServerList = data.List

	app.Node.Register.Set(conn.FD(), register)

	// add to notify queue
	for i := 0; i < len(data.List); i++ {
		app.Node.Alive.AddConn(data.List[i], conn)
	}

	// notify what you are watching
	for i := 0; i < len(data.List); i++ {
		var list = app.Node.Alive.GetData(data.List[i])
		if len(list) == 0 {
			continue
		}

		var err = conn.ProtoBufEmit("/Alive", &message.ServerInfoList{List: list})
		if err != nil {
			console.Error(err)
		}
	}

	return nil
}

func Key(stream *socket.Stream[server.Conn]) error {

	app.Node.Lock()
	defer app.Node.Unlock()

	var conn = stream.Conn

	var data message.KeyList

	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	if len(data.List) == 0 {
		return errors.New("listen list is empty")
	}

	var register = app.Node.Register.Get(conn.FD())
	if register == nil {
		register = &structs.Register{}
	}
	register.KeyList = data.List

	app.Node.Register.Set(conn.FD(), register)

	// add to watch queue
	for i := 0; i < len(data.List); i++ {

		var key = data.List[i]
		app.Node.Key.Add(key, conn)

		var value, err = app.Node.Store.Get(key)
		if err != nil {
			continue
		}

		if value == "" {
			continue
		}

		err = conn.Emit("/Key", []byte(key+"\n"+value))
		if err != nil {
			console.Error(err)
		}
	}

	return nil
}
