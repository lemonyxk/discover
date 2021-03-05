/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-02 15:13
**/

package discover

import (
	"os"
	"time"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/promise"
	"github.com/lemoyxk/utils"

	"discover/app"
	"discover/http"
	"discover/tcp"
	"discover/udp/client"
	"discover/udp/server"
)

func Start(config *app.Config) {

	app.Node.InitConfig(config)

	app.Node.InitRegister()

	app.Node.InitAlive()

	app.Node.InitListen()

	app.Node.InitServerMap()

	app.Node.InitAddr()

	app.Node.InitStore()

	// in this moment, client may talk to you.
	// if you want to do better,
	// you should answer after http server run up.
	// BUT i don't allow save empty value, (sounds good)
	// so is ok.

	// FSM run over and state is not Candidate
	// MAKE SURE that raft server run first
	// or the whole system will FUCK UP
	<-app.Node.Store.Ready

	// start udp server
	// for find others
	var p1 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		server.Start(app.Node.Addr.Udp, func() {
			resolve(nil)
		})
	})

	// http server
	// for set get and delete
	var p2 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		http.Start(app.Node.Addr.Http, func() {
			resolve(nil)
		})
	})

	// udp client
	// send message to others
	var p3 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		client.Start(app.Node.Addr.Udp, func() {
			resolve(nil)
		})
	})

	var p4 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {

		var serverList = app.Node.GetServerList()

		// already finish vote
		if len(serverList) > 0 {
			resolve(nil)
			return
		}

		// one second to decide
		// if it's first tome, that this is very important
		// or that is unnecessary
		client.SendWhoIsMaster()
		var master = app.Node.ServerMap.GetMaster()
		app.Node.Store.BootstrapCluster(master.Addr.Addr == app.Node.Addr.Addr)
		app.Node.Join(master.Addr.Http, app.Node.Addr.Http)

		resolve(nil)
	})

	// waiting master is ready
	var p5 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		for {
			if !app.Node.IsReady() {
				time.Sleep(time.Millisecond * 100)
			}
			resolve(nil)
			break
		}
	})

	// tcp server
	// for consumers and users
	var p6 = promise.New(func(resolve promise.Resolve, reject promise.Reject) {
		tcp.Start(app.Node.Addr.Tcp, func() {
			resolve(nil)
		})
	})

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second)
	// 		log.Println(app.Node.Alive.AllConn())
	// 		log.Println(app.Node.Alive.AllData())
	// 		log.Println(app.Node.Register.All())
	// 		log.Println(app.Node.Listen.All())
	// 	}
	// }()

	promise.Fall(p1, p2, p3, p4, p5, p6).Then(func(result promise.Result) {
		console.Debug("raft server start success")
	})

	utils.Signal.ListenKill().Done(func(sig os.Signal) {
		console.Info("exit with code", sig)
	})
}
