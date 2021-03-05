/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-27 22:07
**/

package discover

import (
	"fmt"
	"strings"
	"time"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/utils"

	"discover/structs"
)

func (dis *discover) getServerList() []structs.WhoIsMaster {

	var rAddr = dis.randomAddr()

	var url = fmt.Sprintf("http://%s/%s", rAddr.Http, "ServerList")

	console.Warning("get server list from", url)

	var res = utils.HttpClient.Get(url).Query(nil).Send()

	if !strings.HasPrefix(res.String(), "OK") {
		if res.LastError() != nil {
			console.Error(res.LastError())
		} else {
			console.Error(res.String()[3:])
		}
		time.Sleep(time.Millisecond * 1000)
		return dis.getServerList()
	}

	var serverList []structs.WhoIsMaster

	_ = utils.Json.Decode(res.Bytes()[3:], &serverList)

	return serverList
}

func (dis *discover) getMasterServer() structs.Address {

	dis.serverList = dis.getServerList()

	var master structs.Address

	for i := 0; i < len(dis.serverList); i++ {
		if dis.serverList[i].IsMaster {
			master = dis.serverList[i].Addr
			break
		}
	}

	if master.Addr == "" {
		time.Sleep(time.Millisecond * 1000)
		return dis.getMasterServer()
	}

	dis.master = master

	return master
}

func (dis *discover) randomAddr() structs.Address {
	var index = utils.Rand.RandomIntn(0, len(dis.serverList))
	return dis.serverList[index].Addr
}

func (dis *discover) url(path string, master bool) string {
	var addr structs.Address

	if master {
		addr = dis.master
	} else {
		addr = dis.randomAddr()
	}

	if addr.Addr == "" {
		panic("random err")
	}

	var url = fmt.Sprintf("http://%s%s", addr.Http, path)
	return url
}
