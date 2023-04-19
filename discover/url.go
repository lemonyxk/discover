/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-27 22:07
**/

package discover

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/socket/http/client"
	"github.com/lemonyxk/utils"
)

// will never be nil
func (dis *discover) getServerList() []*message.Address {
	var rAddr = dis.randomAddr()
	var url = fmt.Sprintf("http://%s/%s", rAddr.Http, "ServerList")

	var res = client.Get(url).Query(nil).Send()
	if res.Error() != nil {
		console.Error(res.Error())
		time.Sleep(time.Millisecond * 1000)
		return dis.getServerList()
	}

	var addr message.AddressResponse
	var err = jsoniter.Unmarshal(res.Bytes(), &addr)
	if err != nil {
		console.Error(err)
		time.Sleep(time.Millisecond * 1000)
		return dis.getServerList()
	}

	if addr.Code != 200 {
		console.Error("get server list error:", addr.Code)
		time.Sleep(time.Millisecond * 1000)
		return dis.getServerList()
	}

	return addr.Msg
}

// will never be nil
func (dis *discover) getMasterServer() *message.Server {

	dis.serverList = dis.getServerList()

	var master *message.Address

	for i := 0; i < len(dis.serverList); i++ {
		if dis.serverList[i].Master {
			master = dis.serverList[i]
			break
		}
	}

	if master == nil {
		time.Sleep(time.Millisecond * 1000)
		return dis.getMasterServer()
	}

	dis.master = master.Server

	return master.Server
}

// will never be nil
func (dis *discover) randomAddr() *message.Server {
	var index = utils.Rand.RandomIntn(0, len(dis.serverList))
	return dis.serverList[index].Server
}

func (dis *discover) url(path string, master bool) string {
	var addr *message.Server

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
