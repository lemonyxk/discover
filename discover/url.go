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
	"github.com/lemonyxk/utils/rand"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/message"
	"github.com/lemonyxk/kitty/socket/http/client"
)

// will never be nil
func (dis *Client) getServerList() []*message.Address {
	var rAddr = dis.randomAddr()
	var url = fmt.Sprintf("http://%s/%s", rAddr.Http, "ServerList")

	var res = client.Get(url).Query().Send()
	if res.Error() != nil {
		console.Error(res.Error())
		time.Sleep(time.Millisecond * 1000)
		return dis.getServerList()
	}

	if res.Code() != 200 {
		console.Error("get server list error:", res.String())
		time.Sleep(time.Millisecond * 1000)
		return dis.getServerList()
	}

	var addr []*message.Address
	var err = jsoniter.Unmarshal(res.Bytes(), &addr)
	if err != nil {
		console.Error(err)
		time.Sleep(time.Millisecond * 1000)
		return dis.getServerList()
	}

	return addr
}

// will never be nil
func (dis *Client) getMasterServer() *message.Server {

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
func (dis *Client) randomAddr() *message.Server {
	var index = rand.RandomIntn(0, len(dis.serverList))
	return dis.serverList[index].Server
}

func (dis *Client) url(path string, master bool) string {
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
