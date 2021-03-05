/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-26 11:15
**/

package client

import (
	"github.com/lemoyxk/kitty/socket/udp/client"
)

func Router(router *client.Router) {
	router.Group().Handler(func(handler *client.RouteHandler) {
		handler.Route("/WhoIsMaster").Handler(WhoIsMaster)
	})
}
