// Package ast @author: Violet-Eva @date  : 2025/5/7 @notes :
package ast

import (
	"strings"
)

type TableName struct {
	catalog string
	db      string
	tbl     string
}

func NewTableNameWithQualifiedName(catalog, database string, qualifiedName QualifiedName) *TableName {
	parts := qualifiedName.GetParts()
	if len(parts) == 3 {
		return NewTableName(parts[0], parts[1], parts[2])
	} else if len(parts) == 2 {
		return NewTableName(catalog, parts[0], parts[1])
	} else if len(parts) == 1 {
		return NewTableName(catalog, database, parts[0])
	}
	return &TableName{}
}

func clean(str string) string {
	if tblCaseSensitive {
		return strings.Trim(str, "`")
	} else {
		return strings.Trim(strings.ToLower(str), "`")
	}
}

func NewTableName(catalog, db, tbl string) *TableName {
	return &TableName{
		catalog: clean(catalog),
		db:      clean(db),
		tbl:     clean(tbl),
	}
}

func (t *TableName) GetTableName() string {
	var sb strings.Builder
	if t.db != "" {
		sb.WriteString(t.db)
		sb.WriteString(".")
	}
	sb.WriteString(t.tbl)
	return sb.String()
}

func (t *TableName) SetCatalog(catalog string) {
	t.catalog = catalog
}

func (t *TableName) GetCatalog() string {
	return t.catalog
}

func (t *TableName) SetDb(db string) {
	t.db = db
}

func (t *TableName) GetDb() string {
	return t.db
}

func (t *TableName) SetTbl(tbl string) {
	t.tbl = tbl
}

func (t *TableName) GetTbl() string {
	return t.tbl
}

func (t *TableName) String() string {
	var sb strings.Builder
	if t.catalog != "" {
		sb.WriteString(t.catalog)
		sb.WriteString(".")
	}
	if t.db != "" {
		sb.WriteString(t.db)
		sb.WriteString(".")
	}
	if t.tbl != "" {
		sb.WriteString(t.tbl)
	}
	return sb.String()
}

func (t *TableName) ToSql() string {
	var sb strings.Builder
	if t.catalog != "" {
		sb.WriteString("`")
		sb.WriteString(t.catalog)
		sb.WriteString("`.")
	}
	if t.db != "" {
		sb.WriteString("`")
		sb.WriteString(t.db)
		sb.WriteString("`.")
	}
	if t.tbl != "" {
		sb.WriteString("`")
		sb.WriteString(t.tbl)
		sb.WriteString("`")
	}
	return sb.String()
}

func (t *TableName) InTableNameList(tblList []TableName) bool {
	for _, tbl := range tblList {
		if strings.EqualFold(t.String(), tbl.String()) {
			return true
		}
	}
	return false
}

func (t *TableName) IsAllDatabase() bool {
	if t.db == "*" {
		return true
	} else {
		return false
	}
}

func (t *TableName) IsAllTable() bool {
	if t.tbl == "*" {
		return true
	} else {
		return false
	}
}

func NewDatabaseWithQualifiedName(catalog string, qualifiedName QualifiedName) *TableName {
	parts := qualifiedName.GetParts()
	if len(parts) == 2 {
		return NewDatabase(parts[0], parts[1])
	} else if len(parts) == 1 {
		return NewDatabase(catalog, parts[0])
	}
	return &TableName{}
}

func NewDatabase(catalog, db string) *TableName {
	return &TableName{
		catalog: clean(catalog),
		db:      clean(db),
	}
}
