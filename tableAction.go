package OrmTS

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
)

func (u *DataUnit) coverPreDefCol(raw tablestore.DefinedColumnSchema) {
	u.key = raw.Name
	typeDesc := raw.ColumnType.ConvertToPbDefinedColumnType
	switch typeDesc() {
	case otsprotocol.DefinedColumnType_DCT_INTEGER:
		u.ValueType = "int"
	case otsprotocol.DefinedColumnType_DCT_BOOLEAN:
		u.ValueType = "bool"
	case otsprotocol.DefinedColumnType_DCT_DOUBLE:
		u.ValueType = "float64"
	case otsprotocol.DefinedColumnType_DCT_STRING:
		u.ValueType = "string"
	case otsprotocol.DefinedColumnType_DCT_BLOB:
		u.ValueType = "[]byte"
	}
}
func (u *DataUnit) coverPk(raw tablestore.PrimaryKeySchema) {
	u.key = *raw.Name
	switch int32(*raw.Type) {
	case int32(otsprotocol.PrimaryKeyType_INTEGER):
		u.ValueType = "int"
	case int32(otsprotocol.PrimaryKeyType_BINARY):
		u.ValueType = "[]byte"
	case int32(otsprotocol.PrimaryKeyType_STRING):
		u.ValueType = "string"

	}
}
