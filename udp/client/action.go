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
	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/udp/client"
	"github.com/lemoyxk/utils"

	"discover/app"
	"discover/structs"
)

func SendWhoIsMaster() {
	for i := 0; i < 10; i++ {
		console.AssertError(
			app.Node.Client.JsonEmit(socket.JsonPack{
				Event: "/WhoIsMaster",
				Data: structs.WhoIsMaster{
					Addr:      *app.Node.Addr,
					Timestamp: app.Node.StartTime.UnixNano(),
					IsMaster:  app.Node.IsMaster(),
				},
			}),
		)

		time.Sleep(time.Millisecond * 100)
	}
}

func WhoIsMaster(c *client.Client, stream *socket.Stream) error {
	var data structs.WhoIsMaster
	var err = utils.Json.Decode(stream.Data, &data)
	if err != nil {
		return err
	}

	app.Node.ServerMap.Set(data.Addr.Addr, data)

	return nil
}
