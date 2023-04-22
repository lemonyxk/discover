/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-21 23:13
**/

package tcp

import (
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

var Middleware = &middleware{}

type middleware struct {
	Controller
}

func (api *middleware) isReady(stream *socket.Stream[server.Conn]) error {
	if !app.Node.IsReady() {
		var msg = "NOT READY"
		return api.Failed(stream.Sender(), stream.Event, msg)
	}
	return nil
}

func (api *middleware) isMaster(stream *socket.Stream[server.Conn]) error {
	if !app.Node.IsMaster() {
		var msg = "NOT MASTER"
		return api.Failed(stream.Sender(), stream.Event, msg)
	}
	return nil
}
