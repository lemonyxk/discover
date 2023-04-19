/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-02 15:13
**/

package discover

import (
	"os"
	"time"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/discover/http"
	"github.com/lemonyxk/discover/tcp"
	"github.com/lemonyxk/promise"
	"github.com/lemonyxk/utils"
)

func Start(config *app.Config) {

	app.Node.InitConfig(config)

	app.Node.InitRegister()

	app.Node.InitAlive()

	app.Node.InitListen()

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

	// http server
	// for set get and delete
	var p2 = promise.New(func(resolve func(any), reject func(error)) {
		http.Start(app.Node.Addr.Http, func() {
			resolve(nil)
		})
	})

	var p4 = promise.New(func(resolve func(any), reject func(error)) {
		var serverList = app.Node.GetServerList()
		// already finish vote
		if len(serverList) > 0 {
			resolve(nil)
			return
		}

		// one second to decide
		// if it's first time, that this is very important
		// or that is unnecessary

		resolve(nil)
	})

	// waiting master is ready
	var p5 = promise.New(func(resolve func(any), reject func(error)) {
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
	var p6 = promise.New(func(resolve func(any), reject func(error)) {
		tcp.Start(app.Node.Addr.Tcp, func() {
			resolve(nil)
		})
	})

	promise.Fall(p2, p4, p5, p6).Then(func(res []any) {
		console.Debug("discover server start success")
	})

	utils.Signal.ListenKill().Done(func(sig os.Signal) {
		console.Info("exit with code", sig)
	})
}
