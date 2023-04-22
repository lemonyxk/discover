/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-19 11:21
**/

package http

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/kitty/socket/http"
)

type Controller struct{}

func (c *Controller) WithCode(stream *http.Stream, code int, msg any) error {
	stream.Response.WriteHeader(code)
	switch v := msg.(type) {
	case []byte:
		return stream.Sender.Bytes(v)
	case string:
		return stream.Sender.String(v)
	}
	var bts, err = jsoniter.Marshal(msg)
	if err != nil {
		return err
	}
	return stream.Sender.Bytes(bts)
}
func (c *Controller) Failed(stream *http.Stream, msg any) error {
	return c.WithCode(stream, 400, msg)
}

func (c *Controller) Error(stream *http.Stream, msg any) error {
	return c.WithCode(stream, 500, msg)
}

func (c *Controller) Forbidden(stream *http.Stream, msg any) error {
	return c.WithCode(stream, 403, msg)
}

func (c *Controller) Success(stream *http.Stream, msg any) error {
	return c.WithCode(stream, 200, msg)
}
