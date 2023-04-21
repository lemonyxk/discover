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
	s.Group().Before(Middleware.localIP, Middleware.secret).Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Get("/IsMaster").Handler(Action.IsMaster)
		handler.Post("/BeMaster").Handler(Action.BeMaster)
		handler.Post("/WhoIsMaster").Handler(Action.WhoIsMaster)
		handler.Get("/Test").Handler(Action.Test)
	})

	s.Group().Before(Middleware.localIP, Middleware.secret, Middleware.isReady, Middleware.isMaster).Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Post("/Joins").Handler(Action.Joins)
		handler.Post("/Join/:addr").Handler(Action.Join)
		handler.Post("/Leave/:addr").Handler(Action.Leave)
	})

	s.Group().Before(Middleware.localIP, Middleware.isReady).Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Get("/ServerList").Handler(Action.ServerList)
		handler.Get("/Get/:key").Handler(Action.Get)
		handler.Get("/All").Handler(Action.All)
	})

	s.Group().Before(Middleware.localIP, Middleware.isReady, Middleware.isMaster).Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Post("/Set/:key").Handler(Action.Set)
		handler.Post("/Delete/:key").Handler(Action.Delete)
	})
}
