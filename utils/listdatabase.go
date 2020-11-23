package utils

import (
	"github.com/yangbojun/pgvacuum/model"
	"database/sql"
	"log"
)

// 定义获取数据库列表函数
func ListDB(db *sql.DB, pgInfo *model.PgConnInfo) []model.PgConnInfo {
	var res []model.PgConnInfo
	var data = pgInfo
	rows, err := db.Query(`SELECT
	datname 
FROM
	pg_database 
WHERE
	datname NOT IN ( 'template0', 'template1' );`)
	defer rows.Close()
	if err != nil {
		log.Fatalf("获取当前数据库实例内所有库失败:%v\n", err)
	}
	for rows.Next() {
		rows.Scan(&data.Database)
		res = append(res, *data)
	}
	return res
}

