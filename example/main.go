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
	"flag"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/discover"
	"github.com/lemonyxk/discover/app"
	"github.com/lemonyxk/utils"
)

func main() {

	var configPath string
	var dir string
	var addr string
	var http string
	var tcp string
	var raft string
	var secret string
	var debug bool

	flag.StringVar(&configPath, "config", "", "config path")
	flag.StringVar(&dir, "dir", "", "data dir")
	flag.StringVar(&addr, "addr", "", "server address")
	flag.StringVar(&http, "http", "", "http address")
	flag.StringVar(&tcp, "tcp", "", "tcp address")
	flag.StringVar(&raft, "raft", "", "raft address")
	flag.StringVar(&secret, "secret", "", "secret key")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	var config app.Config

	if configPath != "" {
		var file = utils.File.ReadFromPath(configPath)
		if file.LastError() != nil {
			console.Exit(file.LastError())
		}
		var err = utils.Json.Decode(file.Bytes(), &config)
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
		config.Debug = debug
	}

	discover.Start(&config)
}
