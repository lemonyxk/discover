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
)

func ServerList(stream *http.Stream) error {
	return stream.Sender.String("OK")
}
