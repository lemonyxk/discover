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
	"log"

	"github.com/lemoyxk/kitty/http"
	"github.com/lemoyxk/utils"

	"discover/app"
)

// if response text is start with "OK"
// that means success
// other means failed

func Join(stream *http.Stream) error {

	var ad = stream.Form.First("addr").String()

	var addr = app.ParseAddr(ad)

	var err = app.Node.Store.Join(addr.Raft)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	return stream.EndString("OK")
}

func Leave(stream *http.Stream) error {

	var ad = stream.Form.First("addr").String()

	var addr = app.ParseAddr(ad)
	log.Println(addr)
	var err = app.Node.Store.Leave(addr.Raft)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	return stream.EndString("OK")
}

func IsMaster(stream *http.Stream) error {
	if app.Node.IsMaster() {
		return stream.EndString("OK")
	}
	return stream.EndString("NO")
}

func WhoIsMaster(stream *http.Stream) error {
	return stream.EndString("OK\n" + app.Node.GetMaster().Addr)
}

// NOTICE
// inaccurate
// because it need time to notify
func ServerList(stream *http.Stream) error {
	return stream.EndString("OK\n" + string(utils.Json.Encode(app.Node.GetServerList())))
}

func Get(stream *http.Stream) error {

	var key = stream.Query.First("key").String()
	if key == "" {
		return stream.EndString("NO\nKEY IS EMPTY")
	}

	value, err := app.Node.Store.Get(key)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	if value == "" {
		return stream.EndString("NO\nVALUE IS EMPTY")
	}

	return stream.EndString("OK\n" + value)
}

func Set(stream *http.Stream) error {

	var key = stream.Form.First("key").String()
	var value = stream.Form.First("value").String()

	if key == "" {
		return stream.EndString("NO\nKEY IS EMPTY")
	}

	if value == "" {
		return stream.EndString("NO\nVALUE IS EMPTY")
	}

	var v, err = app.Node.Store.Get(key)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	if v == value {
		return stream.EndString("OK")
	}

	err = app.Node.Store.Set(key, value)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	return stream.EndString("OK")
}

func Delete(stream *http.Stream) error {
	var key = stream.Form.First("key").String()
	if key == "" {
		return stream.EndString("NO\nKEY IS EMPTY")
	}

	var err = app.Node.Store.Delete(key)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	return stream.EndString("OK")
}

//
// func Login(stream *http.Stream) error {
//
// 	var serverName = stream.Form.First("serverName").String()
// 	if serverName == "" {
// 		return stream.EndString("NO\nSERVER NAME IS EMPTY")
// 	}
//
// 	var addr = stream.Form.First("addr").String()
// 	if addr == "" {
// 		return stream.EndString("NO\nADDR IS EMPTY")
// 	}
//
// 	var err = app.Node.State.Add(serverName, addr)
// 	if err != nil {
// 		return stream.EndString("NO\n" + err.Error())
// 	}
//
// 	return stream.EndString("OK")
// }
//
// func LoginOut(stream *http.Stream) error {
// 	var serverName = stream.Form.First("serverName").String()
// 	if serverName == "" {
// 		return stream.EndString("NO\nSERVER NAME IS EMPTY")
// 	}
//
// 	var addr = stream.Form.First("addr").String()
// 	if addr == "" {
// 		return stream.EndString("NO\nADDR IS EMPTY")
// 	}
//
// 	var err = app.Node.State.Delete(serverName, addr)
// 	if err != nil {
// 		return stream.EndString("NO\n" + err.Error())
// 	}
//
// 	return stream.EndString("OK")
// }
