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

	"github.com/lemoyxk/discover-admin/http"
)

func main() {

	var addr string

	flag.StringVar(&addr, "addr", "", "server address")
	flag.Parse()

	http.Start(addr, func() {

	})
}
