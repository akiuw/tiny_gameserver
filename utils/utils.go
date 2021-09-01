package utils

import (
	"encoding/json"
	. "gameserver/serverlog"
	"os"
)

type serverConfig struct {
	ListenPort 		int		`json:"listen_port"`
	Protocol 		string  `json:"protocol"`
	RoomMaxPlayer 	int     `json:"room_max_player"`
	RedisAddr		string  `json:"redis_addr"`
	RedisPassword	string  `json:"redis_password"`
	LogPath			string  `json:"log_path"`
}

var GServerConfig *serverConfig

func ReadJSONFile(path string, v interface{}) error {
	filePtr, err := os.Open(path)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	GServerConfig = new(serverConfig)
	err := ReadJSONFile("./server_configure.json", GServerConfig)
	if err != nil{
		Error.Println("server_configure.json file error:", err)
		os.Exit(0)
	}
}