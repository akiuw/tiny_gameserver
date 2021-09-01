package serverlog

import (
	"fmt"
	"gameserver/utils"
	"io"
	"log"
	"os"
	"time"
)

var (
	Info *log.Logger
	Warning *log.Logger
	Error * log.Logger
)

func init()  {
	log.SetFlags(log.Ldate|log.Lshortfile)
	nowTime := fmt.Sprintf("%d-%d-%d",time.Now().Month(),time.Now().Day(),time.Now().Hour())
	logfile,err := os.OpenFile(fmt.Sprintf("%s/%s",utils.GServerConfig.LogPath, nowTime),os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err!=nil{
		fmt.Println("日志文件打开失败",err)
	}

	Info = log.New(io.MultiWriter(os.Stdout,logfile),"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stdout,logfile),"Warning:",log.Ldate | log.Ltime | log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr,logfile),"Error:",log.Ldate | log.Ltime | log.Lshortfile)
}
