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
)

func Router(s *router.Router[*http.Stream]) {
	s.Group().Before().Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Get("/ServerList").Handler(ServerList)
	})
}
