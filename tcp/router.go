/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-04 19:11
**/

package tcp

import (
	"github.com/lemoyxk/kitty/socket"
	"github.com/lemoyxk/kitty/socket/websocket/server"

	"discover/app"
)

func Router(router *server.Router) {
	router.Group().Before(isMaster).Handler(func(handler *server.RouteHandler) {
		handler.Route("/Register").Handler(Register)
		handler.Route("/OnRegister").Handler(OnRegister)
	})

	router.Group().Before(isReady).Handler(func(handler *server.RouteHandler) {
		handler.Route("/Listen").Handler(Listen)
	})
}

func isReady(conn *server.Conn, stream *socket.Stream) error {
	if !app.Node.IsReady() {
		var msg = "NO\nNOT READY"
		return conn.JsonEmit(socket.JsonPack{
			Event: stream.Event,
			Data:  msg,
		})
	}
	return nil
}

func isMaster(conn *server.Conn, stream *socket.Stream) error {
	if !app.Node.IsMaster() {
		var msg = "NO\nNOT MASTER"
		return conn.JsonEmit(socket.JsonPack{
			Event: stream.Event,
			Data:  msg,
		})
	}
	return nil
}
