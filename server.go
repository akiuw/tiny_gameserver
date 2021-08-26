package main

import (
	"fmt"
	"gameserver/handles"
	"gameserver/network"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

var g_ServerConfig = make(map[string]interface{})

var Rdb *redis.Client

func init() {
	g_ServerConfig["ip"] = "127.0.0.1"
	g_ServerConfig["port"] = "9000"
	g_ServerConfig["protocol"] = "tcp"
	g_ServerConfig["room_max_player"] = 2

	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := Rdb.Ping().Result()

	if err != nil {
		fmt.Println("error : cant connect to redis server ", err)
		os.Exit(0)
	} else {
		fmt.Println("Info : Redis connect ok!")
	}
}

func StartGame() {

	gs := network.GetInstance()
	gs.SetTimer(time.Microsecond * 10)
	gs.PoolInit(goroutine.Default())
	gs.RegisterProtocol(handles.GetInstance())

	fmt.Println("server init ok")

	log.Fatal(gnet.Serve(gs, "tcp://:9000", gnet.WithMulticore(true), gnet.WithTicker(true)))
}
