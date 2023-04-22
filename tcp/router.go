/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-04 19:11
**/

package tcp

import (
	"github.com/lemonyxk/kitty/router"
	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

func Router(s *router.Router[*socket.Stream[server.Conn]]) {
	s.Group().Before(Middleware.isMaster).Handler(func(handler *router.Handler[*socket.Stream[server.Conn]]) {
		handler.Route("/Register").Handler(Action.Register)
		handler.Route("/Alive").Handler(Action.Alive)
	})

	s.Group().Before(Middleware.isReady).Handler(func(handler *router.Handler[*socket.Stream[server.Conn]]) {
		handler.Route("/Key").Handler(Action.Key)
	})
}
