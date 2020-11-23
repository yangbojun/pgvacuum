package model

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)
//定义数据库连接结构体
type PgConnInfo struct {
	Host string
	Port int
	User string
	Password string
	Database string
}

// 工厂函数
func NewVacuumServer() (*PgConnInfo, int, string) {
	var conn PgConnInfo
	flag.StringVar(&conn.Host, "h", "127.0.0.1", "数据库host")
	flag.IntVar(&conn.Port, "p", 6543, "数据库端口")
	flag.StringVar(&conn.User, "U", "sa", "用户名")
	flag.StringVar(&conn.Password, "W", "tusc@6789#JKL", "密码")
	flag.StringVar(&conn.Database, "d", "postgres", "创建连接时数据库名称")
	var numWorkers int
	flag.IntVar(&numWorkers, "j", 10, "执行vacuum最大并发数")
	if numWorkers < 1 {
		log.Fatalf("检测到并发数为%v不符合规则,请至少开启一个并发数完成vacuum\n")
	}
	var vacuumMod string
	flag.StringVar(&vacuumMod, "m", "analyze",
		"vacuum模式,可以选择analyze(可以简写为a)或者full(可以简写为f)")
	switch vacuumMod {
	case "analyze", "a":
		vacuumMod = "analyze"
	case "full", "f":
		vacuumMod = "full"
	default:
		log.Fatalf("vacuum模式异常,请使用--help查看有哪些模式")
	}
	flag.Parse()
	return &conn, numWorkers, vacuumMod
}

// 连接信息输出函数
func (conn *PgConnInfo) String() string {
	return fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		conn.User, conn.Password, conn.Host, conn.Port, conn.Database)
}

// 获取数据库连接并检查是否可使用
func (conn *PgConnInfo) GetConn() *sql.DB {
	db, err := sql.Open("postgres", conn.String())
	if err != nil {
		log.Fatalf("无法打开数据库连接:%v\n", err)
	}
	//通过db.Ping检查数据库连接地址
	err = db.Ping()
	if err != nil {
		log.Fatalf("%v:%v/%v无法连接，请传入正确的参数:%v\n", conn.Host, conn.Port, conn.Database, err)
	}
	log.Printf("成功连接到数据库,连接地址为:%v:%v/%v\n", conn.Host, conn.Port, conn.Database)
	return db
}