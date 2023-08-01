/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-02 18:01
**/

package http

import (
	"github.com/lemonyxk/kitty/router"
	"github.com/lemonyxk/kitty/socket/http"
	"github.com/lemonyxk/kitty/socket/http/server"
)

func Router(s *router.Router[*http.Stream[server.Conn]]) {
	s.Group().Before().Handler(func(handler *router.Handler[*http.Stream[server.Conn]]) {
		handler.Get("/ServerList").Handler(ServerList)
	})
}
