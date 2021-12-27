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
	"github.com/lemoyxk/discover-admin/app"
	"github.com/lemoyxk/discover/message"
	"github.com/lemoyxk/kitty/http"
	"google.golang.org/protobuf/proto"
)

func UpdateStatus(stream *http.Stream) error {

	var bts = stream.Form.First("data").Bytes()

	var data message.WhoIsMaster
	var err = proto.Unmarshal(bts, &data)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	app.App.ServerMap.Set(data.Addr.Addr, &data)

	return stream.EndString("OK")
}

func WhoIsMaster(stream *http.Stream) error {

	var data = app.App.ServerMap.GetMaster()
	var bts, err = proto.Marshal(data)
	if err != nil {
		return stream.EndString("NO\n" + err.Error())
	}

	return stream.EndString("OK\n" + string(bts))
}
