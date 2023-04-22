/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-19 15:16
**/

package tcp

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/kitty/socket"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

type Controller struct{}

func (c *Controller) WithCode(sender socket.Emitter[server.Conn], event string, code uint32, msg any) error {
	sender.SetCode(code)
	println(sender.Code())
	switch v := msg.(type) {
	case []byte:
		return sender.Emit(event, v)
	case string:
		return sender.Emit(event, []byte(v))
	}
	var bts, err = jsoniter.Marshal(msg)
	if err != nil {
		return err
	}
	return sender.Emit(event, bts)
}

func (c *Controller) Failed(sender socket.Emitter[server.Conn], event string, msg any) error {
	return c.WithCode(sender, event, 400, msg)
}

func (c *Controller) Error(sender socket.Emitter[server.Conn], event string, msg any) error {
	return c.WithCode(sender, event, 500, msg)
}

func (c *Controller) Forbidden(sender socket.Emitter[server.Conn], event string, msg any) error {
	return c.WithCode(sender, event, 403, msg)
}

func (c *Controller) Success(sender socket.Emitter[server.Conn], event string, msg any) error {
	return c.WithCode(sender, event, 200, msg)
}
