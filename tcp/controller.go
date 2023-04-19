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
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/socket/websocket/server"
)

type Controller struct{}

func (c *Controller) WithCode(conn server.Conn, event string, code int, msg any) error {
	return conn.JsonEmit(event, message.Format{Status: "SUCCESS", Code: code, Msg: msg})
}

func (c *Controller) Failed(conn server.Conn, event string, msg any) error {
	return conn.JsonEmit(event, message.Format{Status: "FAILED", Code: 400, Msg: msg})
}

func (c *Controller) Error(conn server.Conn, event string, msg any) error {
	return conn.JsonEmit(event, message.Format{Status: "ERROR", Code: 500, Msg: msg})
}

func (c *Controller) Forbidden(conn server.Conn, event string, msg any) error {
	return conn.JsonEmit(event, message.Format{Status: "FORBIDDEN", Code: 403, Msg: msg})
}

func (c *Controller) Success(conn server.Conn, event string, msg any) error {
	return conn.JsonEmit(event, message.Format{Status: "SUCCESS", Code: 200, Msg: msg})
}
