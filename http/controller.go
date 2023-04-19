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
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/kitty/header"
	"github.com/lemonyxk/kitty/socket/http"
)

type Controller struct{}

func (c *Controller) WithCode(stream *http.Stream, code int, msg any) error {
	return stream.Sender.Json(message.Format{Status: "SUCCESS", Code: code, Msg: msg})
}

func (c *Controller) JsonPretty(stream *http.Stream, msg any) error {
	var format = message.Format{Status: "SUCCESS", Code: 200, Msg: msg}
	stream.Response.Header().Set(header.ContentType, header.ApplicationJson)
	bts, err := jsoniter.MarshalIndent(format, "", "    ")
	if err != nil {
		return err
	}
	_, err = stream.Response.Write(bts)
	return err
}

func (c *Controller) Failed(stream *http.Stream, msg any) error {
	return stream.Sender.Json(message.Format{Status: "FAILED", Code: 400, Msg: msg})
}

func (c *Controller) Error(stream *http.Stream, msg any) error {
	return stream.Sender.Json(message.Format{Status: "ERROR", Code: 500, Msg: msg})
}

func (c *Controller) Forbidden(stream *http.Stream, msg any) error {
	return stream.Sender.Json(message.Format{Status: "FORBIDDEN", Code: 403, Msg: msg})
}

func (c *Controller) Success(stream *http.Stream, msg any) error {
	return stream.Sender.Json(message.Format{Status: "SUCCESS", Code: 200, Msg: msg})
}
