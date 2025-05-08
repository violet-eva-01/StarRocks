// Package ast @author: Violet-Eva @date  : 2025/5/7 @notes :
package ast

import (
	"strings"
)

type Table struct {
	catalog string
	tbl     string
	db      string
	pos     NodePosition
}

func NewTableWithQualifiedName(catalog, database string, qualifiedName QualifiedName) *Table {
	parts := qualifiedName.GetParts()
	if len(parts) == 3 {
		return NewTable(parts[0], parts[1], parts[2], qualifiedName.GetPos())
	} else if len(parts) == 2 {
		return NewTable(catalog, qualifiedName.GetParts()[0], qualifiedName.GetParts()[1],
			qualifiedName.GetPos())
	} else if len(parts) == 1 {
		return NewTable(catalog, database, qualifiedName.GetParts()[0], qualifiedName.GetPos())
	}
	return &Table{}
}

func clean(str string) string {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, "`", "")
	return str
}

func NewTable(catalog, db, tbl string, pos NodePosition) *Table {
	return &Table{
		catalog: clean(catalog),
		db:      clean(db),
		tbl:     clean(tbl),

		pos: pos,
	}
}

func (t *Table) GetDb() string {
	return t.db
}

func (t *Table) GetTbl() string {
	return t.tbl
}

func (t *Table) GetCatalog() string {
	return t.catalog
}

func (t *Table) GetPos() NodePosition {
	return t.pos
}

func (t *Table) String() string {
	var sb strings.Builder
	if t.catalog != "" {
		sb.WriteString(t.catalog)
		sb.WriteString(".")
	}
	if t.db != "" {
		sb.WriteString(t.db)
		sb.WriteString(".")
	}
	sb.WriteString(t.tbl)
	return sb.String()
}

func (t *Table) ToSql() string {
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
	sb.WriteString("`")
	sb.WriteString(t.tbl)
	sb.WriteString("`")
	return sb.String()
}
