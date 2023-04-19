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
	"os"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/utils"
)

func main() {

	var configPath = utils.Args.Get("-f", "--config")
	var dir = utils.Args.Get("--dir")
	var addr = utils.Args.Get("--addr")
	var http = utils.Args.Get("--http")
	var tcp = utils.Args.Get("--tcp")
	var raft = utils.Args.Get("--raft")
	var secret = utils.Args.Get("--secret")

	var config app.Config

	if configPath != "" {
		f, err := os.Open(configPath)
		if err != nil {
			console.Exit(err)
		}
		var file = utils.File.ReadFromReader(f)
		if file.Error() != nil {
			console.Exit(file)
		}
		err = utils.Json.Decode(file.Bytes(), &config)
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

	discover.Start(&config)
}
