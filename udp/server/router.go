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
	"github.com/lemoyxk/kitty/socket/udp/server"
)

func Router(router *server.Router) {
	router.Group().Handler(func(handler *server.RouteHandler) {
		handler.Route("/WhoIsMaster").Handler(WhoIsMaster)
	})
}
