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
	"bytes"
	"fmt"
	"github.com/lemonyxk/kitty/socket/http/server"
	"io"
	"net"

	json "github.com/bytedance/sonic"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/store"
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

func (api *action) Joins(stream *http.Stream[server.Conn]) error {

	var addr []string

	all, err := io.ReadAll(stream.Request.Body)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	err = json.Unmarshal(all, &addr)
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

func (api *action) Join(stream *http.Stream[server.Conn]) error {

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

func (api *action) Leave(stream *http.Stream[server.Conn]) error {

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

func (api *action) IsMaster(stream *http.Stream[server.Conn]) error {
	if app.Node.IsMaster() {
		return api.Success(stream, "OK")
	}
	return api.Success(stream, "NO")
}

func (api *action) WhoIsMaster(stream *http.Stream[server.Conn]) error {
	var master = app.Node.GetMaster()
	return api.Success(stream, master)
}

func (api *action) BeMaster(stream *http.Stream[server.Conn]) error {
	if !app.Node.IsReady() {
		app.Node.Store.BootstrapCluster(true)
		return api.Success(stream, "OK")
	}
	return api.Failed(stream, "ALREADY HAVE MASTER")
}

// ServerList NOTICE
// inaccurate
// because it need time to notify
func (api *action) ServerList(stream *http.Stream[server.Conn]) error {
	var list = app.Node.GetServerList()
	return api.Success(stream, list)
}

func (api *action) Test(stream *http.Stream[server.Conn]) error {
	return api.Success(stream, "OK")
}

func (api *action) Get(stream *http.Stream[server.Conn]) error {

	var key = stream.Params.Get("key")
	if key == "" {
		return api.Failed(stream, "PARAMS ERROR")
	}

	value, err := app.Node.Store.Get(key)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	if len(value) == 0 {
		return api.Failed(stream, "VALUE IS EMPTY")
	}

	return api.Success(stream, value)
}

func (api *action) All(stream *http.Stream[server.Conn]) error {

	var list = app.Node.Store.All()

	var buf = new(bytes.Buffer)

	for i := 0; i < len(list); i++ {
		buf.WriteString(list[i].Key)
		buf.WriteString("\t")
		buf.Write(list[i].Value)
		buf.WriteString("\r\n")
	}

	return api.Success(stream, buf.Bytes())
}

func (api *action) Set(stream *http.Stream[server.Conn]) error {

	var key = stream.Params.Get("key")
	if key == "" {
		return api.Failed(stream, "PARAMS ERROR")
	}

	var value, err = io.ReadAll(stream.Request.Body)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	if len(value) == 0 {
		return api.Failed(stream, "VALUE IS EMPTY")
	}

	v, err := app.Node.Store.Get(key)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	// if value is same, return
	if bytes.Equal(v, value) {
		return api.Success(stream, "OK")
	}

	err = app.Node.Store.Set(key, value)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	return api.Success(stream, "OK")
}

func (api *action) SetMulti(stream *http.Stream[server.Conn]) error {

	var value, err = io.ReadAll(stream.Request.Body)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	if len(value) == 0 {
		return api.Failed(stream, "VALUE IS EMPTY")
	}

	var res []store.KV
	err = json.Unmarshal(value, &res)
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	for i := 0; i < len(res); i++ {
		err = app.Node.Store.Set(res[i].Key, res[i].Value)
		if err != nil {
			return api.Failed(stream, err.Error())
		}
	}

	return api.Success(stream, "OK")
}

func (api *action) Delete(stream *http.Stream[server.Conn]) error {
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

func (api *action) Clear(stream *http.Stream[server.Conn]) error {
	var err = app.Node.Store.Clear()
	if err != nil {
		return api.Failed(stream, err.Error())
	}

	return api.Success(stream, "OK")
}
