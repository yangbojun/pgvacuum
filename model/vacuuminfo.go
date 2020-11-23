package model

import "fmt"

// 定义vacuum语句信息
type VacuumInfo struct {
	PgConnInfo
	VacuumSql string
}

// 定义标准化数据方法
func (this VacuumInfo) String() string {
	return fmt.Sprintf("在%v内执行%v\n", this.Database, this.VacuumSql)
}