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

	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/kitty/v2/router"
	"github.com/lemonyxk/kitty/v2/socket/http"
	"github.com/lemonyxk/utils/v3"
)

func Router(s *router.Router[*http.Stream]) {
	s.Group().Before(localIP, secret).Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Get("/IsMaster").Handler(IsMaster)
		handler.Post("/BeMaster").Handler(BeMaster)
		handler.Get("/Test").Handler(Test)
	})

	s.Group().Before(localIP, secret, ready, isMaster).Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Post("/Join").Handler(Join)
		handler.Post("/Leave").Handler(Leave)
	})

	s.Group().Before(localIP, ready).Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Get("/ServerList").Handler(ServerList)
		handler.Get("/Get").Handler(Get)
	})

	s.Group().Before(localIP, ready, isMaster).Handler(func(handler *router.Handler[*http.Stream]) {
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
