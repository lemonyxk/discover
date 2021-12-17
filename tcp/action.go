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

	"github.com/golang/protobuf/proto"
	"github.com/lemoyxk/console"
	"github.com/lemoyxk/discover/app"
	"github.com/lemoyxk/discover/message"
	"github.com/lemoyxk/discover/structs"
	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/websocket/server"
)

func Register(conn *server.Conn, stream *socket.Stream) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var data message.ServerInfo

	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	if data.ServerName == "" || data.Addr == "" {
		return errors.New("server name or addr is empty")
	}

	var register = app.Node.Register.Get(conn.FD)
	if register == nil {
		register = &structs.Register{}
	}
	register.ServerInfo = &data

	app.Node.Register.Set(conn.FD, register)

	// add to watch queue
	app.Node.Alive.AddData(data.ServerName, data.Addr)
	var list = app.Node.Alive.GetData(data.ServerName)

	var connections = app.Node.Alive.GetConn(data.ServerName)
	for i := 0; i < len(connections); i++ {
		var err = connections[i].ProtoBufEmit(socket.ProtoBufPack{
			Event: "/OnRegister",
			Data:  &message.ServerInfoList{List: list},
		})
		if err != nil {
			console.Error(err)
		}
	}

	return conn.Emit(socket.Pack{
		Event: stream.Event,
		Data:  []byte("OK"),
	})
}

func OnRegister(conn *server.Conn, stream *socket.Stream) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var data message.ServerList

	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	if len(data.List) == 0 {
		return errors.New("server list is empty")
	}

	var register = app.Node.Register.Get(conn.FD)
	if register == nil {
		register = &structs.Register{}
	}
	register.ServerList = data.List

	app.Node.Register.Set(conn.FD, register)

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

		var err = conn.ProtoBufEmit(socket.ProtoBufPack{
			Event: "/OnRegister",
			Data:  &message.ServerInfoList{List: list},
		})
		if err != nil {
			console.Error(err)
		}
	}

	return nil
}

func Listen(conn *server.Conn, stream *socket.Stream) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var data message.KeyList

	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	if len(data.List) == 0 {
		return errors.New("listen list is empty")
	}

	var register = app.Node.Register.Get(conn.FD)
	if register == nil {
		register = &structs.Register{}
	}
	register.KeyList = data.List

	app.Node.Register.Set(conn.FD, register)

	// add to watch queue
	for i := 0; i < len(data.List); i++ {

		var key = data.List[i]
		app.Node.Listen.Add(key, conn)

		var value, err = app.Node.Store.Get(key)
		if err != nil {
			continue
		}

		if value == "" {
			continue
		}

		err = conn.Emit(socket.Pack{
			Event: "/OnListen",
			Data:  []byte(key + "\n" + value),
		})
		if err != nil {
			console.Error(err)
		}
	}

	return conn.Emit(socket.Pack{
		Event: stream.Event,
		Data:  []byte("OK"),
	})
}
