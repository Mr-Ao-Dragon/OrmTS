package OrmTS

import (
	"fmt"
)

func (r Query) ListDatabases() []string {
	dbs := make([]string, 0)
	for _, v := range r.region.Dbs {
		dbs = append(dbs, *v.InstanceName)
	}
	return dbs
}
func (r Query) dbIsExist(db string) bool {

	for _, v := range r.region.Dbs {
		if *v.InstanceName == db {
			return true
		}
	}
	return false
}

// getDbIndexNum 获取目标数据库的索引
// 必须使用上文函数重的判断存在后再进行查询
func (r Query) getDbIndexNum(db string) int {
	for k, v := range r.region.Dbs {
		if *v.InstanceName == db {
			return k
		}
	}
	return -1
}
func (r *Query) ListTables(db string) error {
	if !r.dbIsExist(db) {
		return fmt.Errorf("database is not Exist")
	}
	//targetDbNum := r.getDbIndexNum(db)
	//client := tablestore.NewClientWithConfig(
	//	"https://" + r.region,
	//)
	return nil
}
