/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-02 18:01
**/

package http

import (
	"github.com/lemoyxk/kitty/http/server"
)

func Router(router *server.Router) {
	router.Group().Before().Handler(func(handler *server.RouteHandler) {
		handler.Get("/ServerList").Handler(ServerList)
	})
}
