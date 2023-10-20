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

func Router(s *router.Router[*http.Stream[server.Conn], any]) {
	s.Group().Before(Middleware.localIP, Middleware.secret).Handler(func(handler *router.Handler[*http.Stream[server.Conn], any]) {
		handler.Get("/IsMaster").Handler(Action.IsMaster)
		handler.Post("/BeMaster").Handler(Action.BeMaster)
		handler.Post("/WhoIsMaster").Handler(Action.WhoIsMaster)
		handler.Get("/Test").Handler(Action.Test)
	})

	s.Group().Before(Middleware.localIP, Middleware.secret, Middleware.isReady, Middleware.isMaster).Handler(func(handler *router.Handler[*http.Stream[server.Conn], any]) {
		handler.Post("/Joins").Handler(Action.Joins)
		handler.Post("/Join/:addr").Handler(Action.Join)
		handler.Post("/Leave/:addr").Handler(Action.Leave)
	})

	s.Group().Before(Middleware.localIP, Middleware.isReady).Handler(func(handler *router.Handler[*http.Stream[server.Conn], any]) {
		handler.Get("/ServerList").Handler(Action.ServerList)
		handler.Get("/Get/:key").Handler(Action.Get)
		handler.Get("/All").Handler(Action.All)
	})

	s.Group().Before(Middleware.localIP, Middleware.isReady, Middleware.isMaster).Handler(func(handler *router.Handler[*http.Stream[server.Conn], any]) {
		handler.Post("/Set/:key").Handler(Action.Set)
		handler.Post("/Delete/:key").Handler(Action.Delete)
		handler.Post("/Clear").Handler(Action.Clear)
		handler.Post("/SetMulti").Handler(Action.SetMulti)
	})
}
