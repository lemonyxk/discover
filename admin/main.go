/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-12-27 20:02
**/

package main

import (
	"flag"
	"os"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/discover-admin/http"
	"github.com/lemoyxk/utils"
)

func main() {

	var addr string

	flag.StringVar(&addr, "addr", "", "server address")
	flag.Parse()

	http.Start(addr, func() {})

	utils.Signal.ListenKill().Done(func(sig os.Signal) {
		console.Info("exit with code", sig)
	})
}
