/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-26 11:57
**/

package client

import (
	"time"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/discover/app"
	"github.com/lemoyxk/discover/message"
	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/udp/client"
	"google.golang.org/protobuf/proto"
)

func SendWhoIsMaster() {
	for i := 0; i < 10; i++ {
		var err = app.Node.Client.ProtoBufEmit(socket.ProtoBufPack{
			Event: "/WhoIsMaster",
			Data: &message.WhoIsMaster{
				Addr:      app.Node.Addr,
				Timestamp: app.Node.StartTime.UnixNano(),
				IsMaster:  app.Node.IsMaster(),
			},
		})

		if err != nil {
			console.Error(err)
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func WhoIsMaster(c *client.Client, stream *socket.Stream) error {
	var data message.WhoIsMaster
	var err = proto.Unmarshal(stream.Data, &data)
	if err != nil {
		return err
	}

	app.Node.ServerMap.Set(data.Addr.Addr, &data)

	return nil
}
