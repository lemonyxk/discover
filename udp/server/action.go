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
	"github.com/golang/protobuf/proto"
	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/udp/server"

	"discover/app"
	"discover/message"
)

func WhoIsMaster(conn *server.Conn, stream *socket.Stream) error {

	var data message.WhoIsMaster
	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	app.Node.ServerMap.Set(data.Addr.Addr, &data)

	return conn.ProtoBufEmit(socket.ProtoBufPack{
		Event: "/WhoIsMaster",
		Data: &message.WhoIsMaster{
			Addr:      app.Node.Addr,
			Timestamp: app.Node.StartTime.UnixNano(),
			IsMaster:  app.Node.IsMaster(),
		},
	})
}
