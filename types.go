package OrmTS

import "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
import otsapi "github.com/alibabacloud-go/tablestore-20201209/client"

type Query struct {
	Tables     []string
	Setting    *tablestore.TableStoreClient
	TargetDB   Db
	region     *RegionData
	token      Token
	PreDefRows []string
}
type Column struct {
	Pks  *[]DataUnit
	Rows *[]DataUnit
}
type Table struct {
	TableName string
	Columns   []Column
}
type RegionData struct {
	Region string
	Dbs    []*otsapi.ListInstancesResponseBodyInstances
}
type Db struct {
	metaData otsapi.ListInstancesResponseBodyInstances
	tables   map[string]table
}
type table struct {
	pks  []*DataUnit
	rows []*DataUnit
}
type Token struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      int
}
type DataUnit struct {
	key       string
	Value     interface{}
	ValueType string
}
type Policy struct {
	Version   string `json:"Version"`
	Statement []struct {
		Effect    string   `json:"Effect"`
		Action    []string `json:"Action"`
		Resource  []string `json:"Resource"`
		Condition struct {
			StringEquals struct {
				OtsRegionId string `json:"ots:RegionId"`
			} `json:"StringEquals"`
		} `json:"Condition"`
	} `json:"Statement"`
}
