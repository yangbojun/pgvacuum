package main

import (
	"log"
	"sync"
	"time"
	"github.com/yangbojun/pgvacuum/model"
	"github.com/yangbojun/pgvacuum/utils"
)

func main() {
	// 开启waitgroup用于处理并发
	var wg sync.WaitGroup

	// 通过服务器传入参数生成新的数据库连接信息
	pgConnInfo, numWorkers, vacuumMod := model.NewVacuumServer()

	// 开启带缓冲任务池
	channel := make(chan model.VacuumInfo, numWorkers)

	// 打开数据库连接并做连通性检查
	conn := pgConnInfo.GetConn()

	// 获取所有db列表
	dbList := utils.ListDB(conn, pgConnInfo)

	// 关闭初始化连接
	conn.Close()

	// 并发在将每个库内vacuum信息放入工作任务池内
	exitCount := 0
	for _, db := range dbList {
		wg.Add(1)
		go utils.SetVacuumWork(db, vacuumMod, &wg, channel, &exitCount)
	}

	// 监听创建任务进程,完成后关闭管道
	wg.Add(1)
	go func(exitCount *int, wg *sync.WaitGroup) {
		defer wg.Done()
		for *exitCount < len(dbList) {
			time.Sleep(time.Second)
		}
		log.Print("结束,关闭任务管道")
		close(channel)
	} (&exitCount, &wg)

	// 开启固定并发线程进行vacuum
	for i := 1 ; i <= numWorkers; i++ {
		go func(th int) {
			wg.Add(1)
			for {
				vacuumInfo, notfinish := <- channel
				if !notfinish {
					break
				}
				log.Printf("正在使用协程-%v执行:%v", th, vacuumInfo)
				utils.Vacuum(vacuumInfo)
			}
			wg.Done()
		} (i)
	}

	// 等待所有线程完成处理
	wg.Wait()

	log.Print("本次vacuum任务完成!")
}
