/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-02 18:14
**/

package http

import (
	"github.com/lemonyxk/kitty/socket/http"
	"github.com/lemonyxk/kitty/socket/http/server"
)

func ServerList(stream *http.Stream[server.Conn]) error {
	return stream.Sender.String("OK")
}
