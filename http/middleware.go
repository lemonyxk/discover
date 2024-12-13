/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-19 11:27
**/

package http

import (
	"errors"
	"github.com/lemonyxk/kitty/kitty/header"
	"github.com/lemonyxk/kitty/socket/http/server"
	"github.com/lemonyxk/utils/address"

	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/kitty/socket/http"
)

var Middleware = &middleware{}

type middleware struct {
	Controller
}

func (api *middleware) secret(stream *http.Stream[server.Conn]) error {
	if app.Node.Config.Secret != "" {
		var secret = stream.Request.Header.Get(header.Authorization)
		if secret != app.Node.Config.Secret {
			var msg = "NO PERMISSION"
			_ = api.Forbidden(stream, msg)
			return errors.New(msg)
		}
	}
	return nil
}

func (api *middleware) isMaster(stream *http.Stream[server.Conn]) error {
	if !app.Node.IsMaster() {
		var msg = "NOT MASTER"
		_ = api.Failed(stream, msg)
		return errors.New(msg)
	}
	return nil
}

func (api *middleware) localIP(stream *http.Stream[server.Conn]) error {
	if !address.IsLocalIP(stream.ClientIP()) {
		var msg = "NOT LOCAL IP"
		_ = api.Failed(stream, msg)
		return errors.New(msg)
	}
	return nil
}

func (api *middleware) isReady(stream *http.Stream[server.Conn]) error {
	if !app.Node.IsReady() {
		var msg = "NOT READY"
		_ = api.Failed(stream, msg)
		return errors.New(msg)
	}
	return nil
}
