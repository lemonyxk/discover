/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-02 15:21
**/

package http

import (
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/kitty/v2"
	"github.com/lemonyxk/kitty/v2/socket/http"
	"github.com/lemonyxk/kitty/v2/socket/http/server"
)

func Start(host string, fn func()) {

	var httpServer = server.Server{Addr: host}

	var router = kitty.NewHttpServerRouter()

	Router(router)

	httpServer.Use(func(next server.Middle) server.Middle {
		return func(stream *http.Stream) {
			stream.AutoParse()
			next(stream)
			console.Debug(stream.Request.URL.Path, stream.String())
		}
	})

	httpServer.OnSuccess = func() {
		console.Info("admin server start success", httpServer.LocalAddr())
		fn()
	}

	go httpServer.SetRouter(router).Start()

}
