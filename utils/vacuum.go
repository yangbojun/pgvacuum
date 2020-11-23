package utils

import (
	"fmt"
	"sync"
	"github.com/yangbojun/pgvacuum/model"
	"database/sql"
	"log"
)

// 定义获取当前库中所有vacuum语句函数
func getVacuum(db *sql.DB, vacuumMod string) []string {
	var res []string
	var sqltext string
	sql := fmt.Sprintf(`SELECT 'vacuum '|| '%v' || ' ' || schemaname || '.' || relname || ';'
	FROM
		pg_stat_user_tables
	WHERE
		n_dead_tup > 0 
		AND ( last_vacuum < ( now( ) :: TIMESTAMP + '-1 day' ) OR last_vacuum IS NULL ) 
	ORDER BY
		n_dead_tup DESC ;`, vacuumMod)
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Fatalf("获取vacuum语句失败:%v\n", err)
	}
	for rows.Next() {
		rows.Scan(&sqltext)
		res = append(res, sqltext)
	}
	return res
}

// 遍历所有db将vacuum任务放到任务池内
func SetVacuumWork(pgInfo model.PgConnInfo, vacuumMod string, wg *sync.WaitGroup, pool chan model.VacuumInfo, exitCount *int) {
	// 完成后关闭并发等待
	defer wg.Done()

	//获取连接
	conn := pgInfo.GetConn()
	defer conn.Close()

	// 生成所有vacuum语句
	vacuumSqls := getVacuum(conn, vacuumMod)
	for _, sql := range vacuumSqls {
		vacuumInfo := model.VacuumInfo{
			pgInfo,
			sql,
			}
		pool <- vacuumInfo
	}
	*exitCount ++
}

// 执行vacuum
func Vacuum(vacuumInfo model.VacuumInfo) {
	// 获取数据库连接
	conn := vacuumInfo.GetConn()
	defer conn.Close()

	// 执行vacuum语句
	_, err := conn.Query(vacuumInfo.VacuumSql)
	if err != nil {
		log.Fatalf("执行vacuum失败:%v\n", vacuumInfo)
	}
}