package utils

import (
	"log"
	"os"
	"time"
)

// 初始化日志输出函数
func init() {
	//先判断是否有logs路径,没有则进行创建
	_, err := os.Stat("logs")
	if os.IsNotExist(err) {
		err = os.Mkdir("logs", 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
	date := time.Now().Format("2006-01-02")
	logFileName := "logs/vacuum" + date + ".log"
	logFile, _ := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime)
}