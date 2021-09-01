package configures

import (
	"fmt"
	"gameserver/utils"
	"os"
)

type serverConfig struct {
	ListenPort    int    `json:"listen_port"`
	Protocol      string `json:"protocol"`
	RoomMaxPlayer int    `json:"room_max_player"`
	RedisAddr     string `json:"redis_addr"`
	RedisPassword string `json:"redis_password"`
	LogPath       string `json:"log_path"`
}

var GServerConfig *serverConfig

func init() {
	GServerConfig = new(serverConfig)
	err := utils.ReadJSONFile("./server_configure.json", GServerConfig)
	if err != nil {
		fmt.Println("server_configure.json open failed!", err)
		os.Exit(0)
	}
}
