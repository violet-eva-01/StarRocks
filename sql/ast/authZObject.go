// Package ast @author: Violet-Eva @date  : 2025/5/23 @notes :
package ast

type Object struct {
	OnType     map[privilegeOnType]PrivilegeObjectTypePlural
	ObjectType PrivilegeObjectType
	DbName     string
	TblName    string
	OtherName  string
}

func NewEmptyObject() *Object {
	return &Object{}
}

func GetObject(tb TableName, otherName string) *Object {
	obj := NewEmptyObject()
	obj.OtherName = otherName
	obj.DbName = tb.db
	obj.TblName = tb.tbl
	return obj
}
