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
	"net"

	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/kitty/socket/http"
	"github.com/lemonyxk/kitty/socket/http/client"
	"github.com/lemonyxk/utils/v3"
)

// if response text is start with "OK"
// that means success
// other means failed

func Join(stream *http.Stream) error {

	var ad = stream.Form.First("addr").String()

	if ad == "" {
		return stream.Sender.String("NO\n" + "addr is empty")
	}

	_, err := net.ResolveTCPAddr("tcp", ad)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	var addr = app.ParseAddr(ad)

	var res = client.Get("http://" + addr.Http + "/Test").Query().Send().String()
	if res != "OK\n" {
		return stream.Sender.String("NO\n" + "addr no response")
	}

	err = app.Node.Store.Join(addr.Raft)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	return stream.Sender.String("OK")
}

func Leave(stream *http.Stream) error {

	var ad = stream.Form.First("addr").String()

	if ad == "" {
		return stream.Sender.String("NO\n" + "addr is empty")
	}

	_, err := net.ResolveTCPAddr("tcp", ad)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	var addr = app.ParseAddr(ad)

	err = app.Node.Store.Leave(addr.Raft)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	return stream.Sender.String("OK")
}

func IsMaster(stream *http.Stream) error {
	if app.Node.IsMaster() {
		return stream.Sender.String("OK")
	}
	return stream.Sender.String("NO")
}

func BeMaster(stream *http.Stream) error {
	if !app.Node.IsReady() {
		app.Node.Store.BootstrapCluster(true)
		return stream.Sender.String("OK")
	}
	return stream.Sender.String("NO")
}

// ServerList NOTICE
// inaccurate
// because it need time to notify
func ServerList(stream *http.Stream) error {
	return stream.Sender.String("OK\n" + string(utils.Json.Encode(app.Node.GetServerList())))
}

func Test(stream *http.Stream) error {
	return stream.Sender.String("OK\n")
}

func Get(stream *http.Stream) error {

	var key = stream.Query.First("key").String()
	if key == "" {
		return stream.Sender.String("NO\nKEY IS EMPTY")
	}

	value, err := app.Node.Store.Get(key)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	if value == "" {
		return stream.Sender.String("NO\nVALUE IS EMPTY")
	}

	return stream.Sender.String("OK\n" + value)
}

func Set(stream *http.Stream) error {

	var key = stream.Form.First("key").String()
	var value = stream.Form.First("value").String()

	if key == "" {
		return stream.Sender.String("NO\nKEY IS EMPTY")
	}

	if value == "" {
		return stream.Sender.String("NO\nVALUE IS EMPTY")
	}

	var v, err = app.Node.Store.Get(key)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	if v == value {
		return stream.Sender.String("OK")
	}

	err = app.Node.Store.Set(key, value)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	return stream.Sender.String("OK")
}

func Delete(stream *http.Stream) error {
	var key = stream.Form.First("key").String()
	if key == "" {
		return stream.Sender.String("NO\nKEY IS EMPTY")
	}

	var err = app.Node.Store.Delete(key)
	if err != nil {
		return stream.Sender.String("NO\n" + err.Error())
	}

	return stream.Sender.String("OK")
}

func All(stream *http.Stream) error {
	var list = app.Node.Store.All()
	return stream.Sender.String("OK\n" + string(utils.Json.Encode(list)))
}
