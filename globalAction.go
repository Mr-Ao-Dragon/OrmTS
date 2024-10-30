package OrmTS

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
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
// 必须使用上文函数中的判断存在后再进行查询
func (r Query) getDbIndexNum(db string) int {
	for k, v := range r.region.Dbs {
		if *v.InstanceName == db {
			return k
		}
	}
	return -1
}
func (r *Query) FetchTables(db string) error {
	if !r.dbIsExist(db) {
		return fmt.Errorf("database is not Exist")
	}
	targetDbNum := r.getDbIndexNum(db)
	// 初始化客户端
	client := tablestore.NewClientWithConfig(
		"https://"+r.region.Region,
		db,
		r.token.AccessKeyId,
		r.token.AccessKeySecret,
		r.token.SecurityToken,
		nil,
	)
	// 获取表名
	tablename, err := client.ListTable()
	if err != nil {
		return errors.Join(err, fmt.Errorf("fetch tables name failed"))
	}
	r.TargetDB.metaData = *r.region.Dbs[targetDbNum]
	describ := new(tablestore.DescribeTableResponse)
	describeTableReq := new(tablestore.DescribeTableRequest)
	// 遍历表名
	for _, v := range tablename.TableNames {
		describeTableReq.TableName = v
		var errsig error
		// 获取本次迭代的表结构
		describ, errsig = client.DescribeTable(describeTableReq)
		if errsig != nil {
			errors.Join(err, errsig)
		}
		// 转换出主键
		for Index := range describ.TableMeta.SchemaEntry {
			r.TargetDB.tables[describ.TableMeta.TableName].pks[Index].coverPk(*describ.TableMeta.SchemaEntry[Index])
		}
		// 转换出预定义列
		for index := range describ.TableMeta.DefinedColumns {
			r.TargetDB.tables[describ.TableMeta.TableName].rows[index].coverPreDefCol(*describ.TableMeta.DefinedColumns[index])
		}
	}
	return nil
}
