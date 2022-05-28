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
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/kitty/v2/router"
	"github.com/lemonyxk/kitty/v2/socket"
	"github.com/lemonyxk/kitty/v2/socket/websocket/server"
)

func Router(s *router.Router[*socket.Stream[server.Conn]]) {
	s.Group().Before(isMaster).Handler(func(handler *router.Handler[*socket.Stream[server.Conn]]) {
		handler.Route("/Register").Handler(Register)
		handler.Route("/Alive").Handler(Alive)
	})

	s.Group().Before(isReady).Handler(func(handler *router.Handler[*socket.Stream[server.Conn]]) {
		handler.Route("/Key").Handler(Key)
	})
}

func isReady(stream *socket.Stream[server.Conn]) error {
	if !app.Node.IsReady() {
		var msg = "NO\nNOT READY"
		return stream.Conn.Emit(stream.Event, []byte(msg))
	}
	return nil
}

func isMaster(stream *socket.Stream[server.Conn]) error {
	if !app.Node.IsMaster() {
		var msg = "NO\nNOT MASTER"
		return stream.Conn.Emit(stream.Event, []byte(msg))
	}
	return nil
}
