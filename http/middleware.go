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

	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/kitty/socket/http"
	"github.com/lemonyxk/utils"
)

var Middleware = &middleware{}

type middleware struct {
	Controller
}

func (api *middleware) secret(stream *http.Stream) error {
	var secret = stream.AutoGet("secret").String()
	if app.Node.Config.Secret != secret {
		var msg = "NO PERMISSION"
		_ = api.Forbidden(stream, msg)
		return errors.New(msg)
	}
	return nil
}

func (api *middleware) isMaster(stream *http.Stream) error {
	if !app.Node.IsMaster() {
		var msg = "NOT MASTER"
		_ = api.Failed(stream, msg)
		return errors.New(msg)
	}
	return nil
}

func (api *middleware) localIP(stream *http.Stream) error {
	if !utils.Addr.IsLocalIP(stream.ClientIP()) {
		var msg = "NOT LOCAL IP"
		_ = api.Failed(stream, msg)
		return errors.New(msg)
	}
	return nil
}

func (api *middleware) isReady(stream *http.Stream) error {
	if !app.Node.IsReady() {
		var msg = "NOT READY"
		_ = api.Failed(stream, msg)
		return errors.New(msg)
	}
	return nil
}
