/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2023-04-19 16:21
**/

package message

type Format struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    any    `json:"msg"`
}

type OpResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    Op     `json:"msg"`
}

type ServerInfoResponse struct {
	Status string        `json:"status"`
	Code   int           `json:"code"`
	Msg    []*ServerInfo `json:"msg"`
}

type AddressResponse struct {
	Status string     `json:"status"`
	Code   int        `json:"code"`
	Msg    []*Address `json:"msg"`
}
