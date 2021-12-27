/**
* @program: discover
*
* @description:
*
* @author: lemo
*
* @create: 2021-02-02 15:24
**/

package main

import (
	"flag"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/discover"
	"github.com/lemoyxk/discover/app"
	"github.com/lemoyxk/utils"
)

func main() {

	var configPath string
	var dir string
	var admin string
	var addr string
	var secret string
	var debug bool

	flag.StringVar(&configPath, "config", "", "config path")
	flag.StringVar(&dir, "dir", "", "data dir")
	flag.StringVar(&admin, "admin", "", "admin address")
	flag.StringVar(&addr, "addr", "", "server address")
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
		config.Admin = admin
		config.Secret = secret
		config.Debug = debug
	}

	discover.Start(&config)
}
