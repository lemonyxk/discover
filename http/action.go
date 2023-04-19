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
	"fmt"
	"io"
	"net"

	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/kitty/socket/http"
	"github.com/lemonyxk/kitty/socket/http/client"
)

// if response text is start with "OK"
// that means success
// other means failed

var Action = new(action)

type action struct {
	Controller
}

func (api *action) Joins(stream *http.Stream) error {

	var addr []string

	all, err := io.ReadAll(stream.Request.Body)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	err = jsoniter.Unmarshal(all, &addr)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	if len(addr) == 0 {
		return api.Failed(stream, "addr is empty")
	}

	var errMap = make(map[string]string)

	for i := 0; i < len(addr); i++ {
		_, err := net.ResolveTCPAddr("tcp", addr[i])
		if err != nil {
			errMap[addr[i]] = fmt.Sprintf("addr %s is invalid: %s", addr[i], err.Error())
			continue
		}

		var ad = app.ParseAddr(addr[i])
		var res = client.Get("http://" + ad.Http + "/Test").Query().Send().String()
		if res != "OK" {
			errMap[addr[i]] = fmt.Sprintf("addr %s is not ready", addr[i])
			continue
		}

		err = app.Node.Store.Join(ad.Raft)
		if err != nil {
			errMap[addr[i]] = fmt.Sprintf("addr %s join failed: %s", addr[i], err.Error())
			continue
		}

		errMap[addr[i]] = "OK"
	}

	return api.Success(stream, errMap)
}

func (api *action) Join(stream *http.Stream) error {

	var addr = stream.Params.Get("addr")
	if addr == "" {
		return api.Failed(stream, "PARAMS ERROR")
	}

	_, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return api.Failed(stream, fmt.Sprintf("addr %s is invalid: %s", addr, err.Error()))
	}

	var ad = app.ParseAddr(addr)
	var res = client.Get("http://" + ad.Http + "/Test").Query().Send().String()
	if res != "OK" {
		return api.Failed(stream, fmt.Sprintf("addr %s is not ready", addr))
	}

	err = app.Node.Store.Join(ad.Raft)
	if err != nil {
		return api.Failed(stream, fmt.Sprintf("addr %s join failed: %s", addr, err.Error()))
	}

	return api.Success(stream, "OK")
}

func (api *action) Leave(stream *http.Stream) error {

	var addr = stream.Params.Get("addr")
	if addr == "" {
		return api.Failed(stream, "PARAMS ERROR")
	}

	_, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	var ad = app.ParseAddr(addr)

	err = app.Node.Store.Leave(ad.Raft)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	return api.Success(stream, "OK")
}

func (api *action) IsMaster(stream *http.Stream) error {
	if app.Node.IsMaster() {
		return api.Success(stream, "OK")
	}
	return api.Success(stream, "NO")
}

func (api *action) BeMaster(stream *http.Stream) error {
	if !app.Node.IsReady() {
		app.Node.Store.BootstrapCluster(true)
		return api.Success(stream, "OK")
	}
	return api.Failed(stream, "ALREADY HAVE MASTER")
}

// ServerList NOTICE
// inaccurate
// because it need time to notify
func (api *action) ServerList(stream *http.Stream) error {
	var list = app.Node.GetServerList()
	return api.JsonPretty(stream, list)
}

func (api *action) Test(stream *http.Stream) error {
	return stream.Sender.String("OK")
}

func (api *action) Get(stream *http.Stream) error {

	var key = stream.Params.Get("key")
	if key == "" {
		return api.Failed(stream, "PARAMS ERROR")
	}

	value, err := app.Node.Store.Get(key)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	if value == "" {
		return api.Failed(stream, "VALUE IS EMPTY")
	}

	return api.Success(stream, value)
}

func (api *action) Set(stream *http.Stream) error {

	var key = stream.Params.Get("key")
	if key == "" {
		return api.Failed(stream, "PARAMS ERROR")
	}

	var all, err = io.ReadAll(stream.Request.Body)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	var value = string(all)

	if value == "" {
		return api.Failed(stream, "VALUE IS EMPTY")
	}

	v, err := app.Node.Store.Get(key)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	if v == value {
		return api.Success(stream, "OK")
	}

	err = app.Node.Store.Set(key, value)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	return api.Success(stream, "OK")
}

func (api *action) Delete(stream *http.Stream) error {
	var key = stream.Params.Get("key")
	if key == "" {
		return api.Failed(stream, "PARAMS ERROR")
	}

	var err = app.Node.Store.Delete(key)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	return api.Success(stream, "OK")
}

func (api *action) All(stream *http.Stream) error {
	var list = app.Node.Store.All()
	return api.JsonPretty(stream, list)
}
