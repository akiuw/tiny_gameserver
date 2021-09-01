package main

import (
	"fmt"
	. "gameserver/configures"
	"gameserver/handles"
	"gameserver/network"
	. "gameserver/serverlog"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

var Rdb *redis.Client

func init() {

	Rdb = redis.NewClient(&redis.Options{
		Addr:     GServerConfig.RedisAddr,
		Password: GServerConfig.RedisPassword, // no password set
		DB:       0,                           // use default DB
	})

	_, err := Rdb.Ping().Result()

	if err != nil {
		Error.Println("error : cant connect to redis server ", err)
		os.Exit(0)
	} else {
		Info.Println("Info : Redis connect ok!")
	}
}

func StartGame() {

	gs := network.GetInstance()
	gs.SetTimer(time.Microsecond * 10)
	gs.PoolInit(goroutine.Default())
	gs.RegisterProtocol(handles.GetInstance())

	Info.Println("server init ok")

	log.Fatal(gnet.Serve(gs, fmt.Sprintf("%s://:%d", GServerConfig.Protocol, GServerConfig.ListenPort),
		gnet.WithMulticore(true), gnet.WithTicker(true)))
}
