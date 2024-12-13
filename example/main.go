/**
* @program: discover
*
* @description:
*
* @author: lemon
*
* @create: 2021-02-02 15:24
**/

package main

import (
	json "github.com/bytedance/sonic"
	"github.com/lemonyxk/utils/args"
	file2 "github.com/lemonyxk/utils/file"
	"os"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover"
	"github.com/lemonyxk/discover/app"
)

func main() {

	var configPath = args.Get("-f", "--config")
	var dir = args.Get("--dir")
	var addr = args.Get("--addr")
	var http = args.Get("--http")
	var tcp = args.Get("--tcp")
	var raft = args.Get("--raft")
	var secret = args.Get("--secret")

	var config app.Config

	if configPath != "" {
		f, err := os.Open(configPath)
		if err != nil {
			console.Exit(err)
		}
		var file = file2.FromReader(f)
		if file.Error() != nil {
			console.Exit(file)
		}
		err = json.Unmarshal(file.Bytes(), &config)
		if err != nil {
			console.Exit(err)
		}
	} else {
		config.Dir = dir
		config.Addr = addr
		config.Http = http
		config.Tcp = tcp
		config.Raft = raft
		config.Secret = secret
	}

	console.DefaultLogger.AddField("addr", config.Addr)

	discover.Start(&config)
}
