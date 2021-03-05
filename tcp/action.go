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

	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/websocket/server"
	"github.com/lemoyxk/utils"

	"discover/app"
	"discover/structs"
)

func Register(conn *server.Conn, stream *socket.Stream) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var data structs.ServerInfo

	var err = utils.Json.Decode(stream.Data, &data)
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
		_ = connections[i].JsonEmit(socket.JsonPack{
			Event: "/OnRegister",
			Data:  list,
		})
	}

	return conn.JsonEmit(socket.JsonPack{
		Event: stream.Event,
		Data:  "OK",
	})
}

func OnRegister(conn *server.Conn, stream *socket.Stream) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var data []string

	var err = utils.Json.Decode(stream.Data, &data)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("server list is empty")
	}

	var register = app.Node.Register.Get(conn.FD)
	if register == nil {
		register = &structs.Register{}
	}
	register.ServerList = data

	app.Node.Register.Set(conn.FD, register)

	// add to notify queue
	for i := 0; i < len(data); i++ {
		app.Node.Alive.AddConn(data[i], conn)
	}

	// notify what you are watching
	for i := 0; i < len(data); i++ {
		var list = app.Node.Alive.GetData(data[i])
		if len(list) == 0 {
			continue
		}

		_ = conn.JsonEmit(socket.JsonPack{
			Event: "/OnRegister",
			Data:  list,
		})
	}

	return nil
}

func Listen(conn *server.Conn, stream *socket.Stream) error {

	app.Node.Lock()

	defer app.Node.Unlock()

	var data []string

	var err = utils.Json.Decode(stream.Data, &data)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("listen list is empty")
	}

	var register = app.Node.Register.Get(conn.FD)
	if register == nil {
		register = &structs.Register{}
	}
	register.KeyList = data

	app.Node.Register.Set(conn.FD, register)

	// add to watch queue
	for i := 0; i < len(data); i++ {
		app.Node.Listen.Add(data[i], conn)

		var value, err = app.Node.Store.Get(data[i])
		if err != nil {
			continue
		}

		if value == "" {
			continue
		}

		_ = conn.Emit(socket.Pack{
			Event: "/OnListen",
			Data:  []byte(value),
		})
	}

	return conn.JsonEmit(socket.JsonPack{
		Event: stream.Event,
		Data:  "OK",
	})
}
