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
	"errors"

	"github.com/lemoyxk/discover/app"
	"github.com/lemoyxk/kitty/http"
	"github.com/lemoyxk/kitty/http/server"
	"github.com/lemoyxk/utils"
)

func Router(router *server.Router) {
	router.Group().Before(localIP, secret).Handler(func(handler *server.RouteHandler) {
		handler.Get("/IsMaster").Handler(IsMaster)
		handler.Post("/BeMaster").Handler(BeMaster)
		handler.Get("/Test").Handler(Test)
	})

	router.Group().Before(localIP, secret, ready, isMaster).Handler(func(handler *server.RouteHandler) {
		handler.Post("/Join").Handler(Join)
		handler.Post("/Leave").Handler(Leave)
	})

	router.Group().Before(localIP, ready).Handler(func(handler *server.RouteHandler) {
		handler.Get("/ServerList").Handler(ServerList)
		handler.Get("/Get").Handler(Get)
	})

	router.Group().Before(localIP, ready, isMaster).Handler(func(handler *server.RouteHandler) {
		handler.Post("/Set").Handler(Set)
		handler.Post("/Delete").Handler(Delete)
	})
}

func secret(stream *http.Stream) error {
	var secret = stream.AutoGet("secret").String()
	if app.Node.Config.Secret != secret {
		var msg = "NO\nNO PERMISSION"
		_ = stream.EndString(msg)
		return errors.New(msg)
	}
	return nil
}

func isMaster(stream *http.Stream) error {
	if !app.Node.IsMaster() {
		var msg = "NO\nNOT MASTER"
		_ = stream.EndString(msg)
		return errors.New(msg)
	}
	return nil
}

func localIP(stream *http.Stream) error {
	if !utils.Addr.IsLocalIP(stream.ClientIP()) {
		var msg = "NO\nNOT LOCAL IP"
		_ = stream.EndString(msg)
		return errors.New(msg)
	}
	return nil
}

func ready(stream *http.Stream) error {
	if !app.Node.IsReady() {
		var msg = "NO\nNOT READY"
		_ = stream.EndString(msg)
		return errors.New(msg)
	}
	return nil
}
