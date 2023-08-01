/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-12-27 20:02
**/

package main

import (
	"flag"
	"github.com/lemonyxk/utils/signal"
	"os"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover-admin/http"
)

func main() {

	var addr string

	flag.StringVar(&addr, "addr", "", "server address")
	flag.Parse()

	http.Start(addr, func() {})

	signal.ListenKill().Done(func(sig os.Signal) {
		console.Info("exit with code", sig)
	})
}
