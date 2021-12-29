/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-02 18:14
**/

package http

import (
	"github.com/lemoyxk/kitty/http"
)

func ServerList(stream *http.Stream) error {
	return stream.EndString("OK")
}
