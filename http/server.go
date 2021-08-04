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
	"github.com/lemoyxk/console"
	"github.com/lemoyxk/kitty/http"
	"github.com/lemoyxk/kitty/http/server"
)

// var transport = http2.Transport{
// 	TLSHandshakeTimeout:   10 * time.Second,
// 	ResponseHeaderTimeout: 15 * time.Second,
// 	ExpectContinueTimeout: 2 * time.Second,
// 	MaxIdleConns:          runtime.NumCPU() * 2,
// 	MaxIdleConnsPerHost:   runtime.NumCPU() * 2,
// 	MaxConnsPerHost:       runtime.NumCPU() * 2,
// }
//
// func proxy(stream *http.Stream) {
// 	var ip, port, err = utils.Addr.Parse(string(node.Node.Raft().Leader()))
// 	exception.AssertError(err)
// 	var host = fmt.Sprintf("%s:%d", ip, port-1000)
// 	var proxy = httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: host})
// 	proxy.Transport = &transport
// 	proxy.ServeHTTP(stream.Response, stream.Request)
// }

func Start(host string, fn func()) {

	var httpServer = server.Server{Addr: host}

	var router = server.Router{IgnoreCase: true}

	Router(&router)

	httpServer.Use(func(next server.Middle) server.Middle {
		return func(stream *http.Stream) {
			// if stream.Request.Method == "POST" {
			// 	if node.Node.Raft().State() != raft.Leader {
			// 		proxy(stream)
			// 		return
			// 	}
			// }
			stream.AutoParse()
			next(stream)
			console.Debug(stream.Request.URL.Path, stream.String())
		}
	})

	httpServer.OnSuccess = func() {
		console.Info("http server start success", httpServer.LocalAddr())
		fn()
	}

	go httpServer.SetRouter(&router).Start()

}
