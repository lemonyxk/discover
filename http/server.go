/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-02 15:21
**/

package http

import (
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/kitty"
	"github.com/lemonyxk/kitty/socket/http"
	"github.com/lemonyxk/kitty/socket/http/server"
)

func Start(host string, fn func()) {

	var httpServer = server.Server{Addr: host}

	var router = kitty.NewHttpServerRouter()

	Router(router)

	httpServer.Use(func(next server.Middle) server.Middle {
		return func(stream *http.Stream) {
			// proxy websocket to another port
			next(stream)
		}
	})

	httpServer.OnSuccess = func() {
		console.Info("http server start success", host)
		fn()
	}

	go httpServer.SetRouter(router).Start()

}
