/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-26 11:15
**/

package server

import (
	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/udp/server"
	"github.com/lemoyxk/utils"

	"discover/app"
	"discover/structs"
)

func WhoIsMaster(conn *server.Conn, stream *socket.Stream) error {

	var data structs.WhoIsMaster
	var err = utils.Json.Decode(stream.Data, &data)
	if err != nil {
		return err
	}

	app.Node.ServerMap.Set(data.Addr.Addr, data)

	return conn.JsonEmit(socket.JsonPack{
		Event: "/WhoIsMaster",
		Data: structs.WhoIsMaster{
			Addr:      *app.Node.Addr,
			Timestamp: app.Node.StartTime.UnixNano(),
			IsMaster:  app.Node.IsMaster(),
		},
	})
}
