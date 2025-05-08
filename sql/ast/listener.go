// Package ast @author: Violet-Eva @date  : 2025/5/8 @notes :
package ast

import (
	"github.com/violet-eva-01/StarRocks/sql/parser"
	"strings"
)

const (
	DefaultCatalog  = "default_catalog"
	DefaultDatabase = "default_db"
)

type Listener struct {
	defaultCatalog  string
	defaultDatabase string
	Tables          []Table
	ActionTables    map[Action][]Table
	*parser.BaseStarRocksListener
}

func (l *Listener) SetCatalog(catalog string) {
	l.defaultCatalog = catalog
}

func (l *Listener) SetDatabase(database string) {
	l.defaultDatabase = database
}

func (l *Listener) GetCatalog() string {
	return l.defaultCatalog
}

func (l *Listener) GetDatabase() string {
	return l.defaultDatabase
}

func (l *Listener) GetTables() []Table {
	return l.Tables
}

func (l *Listener) GetActionTables() map[Action][]Table {
	return l.ActionTables
}

func (l *Listener) GetTableNamesByAction(action Action) []Table {
	return l.ActionTables[action]
}

func NewListener() *Listener {
	return &Listener{
		defaultCatalog:  DefaultCatalog,
		defaultDatabase: DefaultDatabase,
		ActionTables:    make(map[Action][]Table),
	}
}

func (l *Listener) getTableNameWithAction(strAction string, ctx parser.IQualifiedNameContext) {
	action := GetActionFromString(strAction)
	tbl := NewTableWithQualifiedName(l.defaultCatalog, l.defaultDatabase, GetQualifiedName(ctx))
	if l.ActionTables[action] == nil {
		l.ActionTables[action] = []Table{*tbl}
	} else {
		l.ActionTables[action] = append(l.ActionTables[action], *tbl)
	}
	l.Tables = append(l.Tables, *tbl)
}

func (l *Listener) EnterSetCatalogStatement(ctx *parser.SetCatalogStatementContext) {
	catalog := strings.ToLower(ctx.IdentifierOrString().Identifier().GetText())
	l.defaultCatalog = catalog
}

func (l *Listener) EnterUseDatabaseStatement(ctx *parser.UseDatabaseStatementContext) {
	parts := GetQualifiedName(ctx.QualifiedName()).GetParts()
	if len(parts) == 2 {
		l.defaultCatalog = strings.ToLower(parts[0])
		l.defaultDatabase = strings.ToLower(parts[1])
	} else if len(parts) == 1 {
		l.defaultDatabase = strings.ToLower(parts[0])
	}
}

// Table DDL Statement

func (l *Listener) EnterCreateTableStatement(ctx *parser.CreateTableStatementContext) {
	l.getTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterCreateTableAsSelectStatement(ctx *parser.CreateTableAsSelectStatementContext) {
	l.getTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterCreateTableLikeStatement(ctx *parser.CreateTableLikeStatementContext) {
	l.getTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName(0))
	l.getTableNameWithAction("SELECT", ctx.QualifiedName(1))
}

func (l *Listener) EnterAlterTableStatement(ctx *parser.AlterTableStatementContext) {
	l.getTableNameWithAction(ctx.ALTER().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterTruncateTableStatement(ctx *parser.TruncateTableStatementContext) {
	l.getTableNameWithAction(ctx.TRUNCATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterDropTableStatement(ctx *parser.DropTableStatementContext) {
	l.getTableNameWithAction(ctx.DROP().GetText(), ctx.QualifiedName())
}

// View Statement

func (l *Listener) EnterCreateViewStatement(ctx *parser.CreateViewStatementContext) {
	l.getTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterDropViewStatement(ctx *parser.DropViewStatementContext) {
	l.getTableNameWithAction(ctx.DROP().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterAlterViewStatement(ctx *parser.AlterViewStatementContext) {
	l.getTableNameWithAction(ctx.ALTER().GetText(), ctx.QualifiedName())
}

// Materialized View Statement

func (l *Listener) EnterCreateMaterializedViewStatement(ctx *parser.CreateMaterializedViewStatementContext) {
	l.getTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterDropMaterializedViewStatement(ctx *parser.DropMaterializedViewStatementContext) {
	l.getTableNameWithAction(ctx.DROP().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterAlterMaterializedViewStatement(ctx *parser.AlterMaterializedViewStatementContext) {
	l.getTableNameWithAction(ctx.ALTER().GetText(), ctx.QualifiedName())
}

// DML Statement

func (l *Listener) EnterTableAtom(ctx *parser.TableAtomContext) {
	l.getTableNameWithAction("SELECT", ctx.QualifiedName())
}

func (l *Listener) EnterInsertStatement(ctx *parser.InsertStatementContext) {
	l.getTableNameWithAction(ctx.INSERT().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterUpdateStatement(ctx *parser.UpdateStatementContext) {
	l.getTableNameWithAction(ctx.UPDATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterDeleteStatement(ctx *parser.DeleteStatementContext) {
	l.getTableNameWithAction(ctx.DELETE().GetText(), ctx.QualifiedName())
}

// AuthZ Statement

func (l *Listener) EnterExecuteAsStatement(ctx *parser.ExecuteAsStatementContext) {}

func (l *Listener) EnterGrantRoleToUser(ctx *parser.GrantRoleToUserContext) {}

func (l *Listener) EnterGrantRoleToRole(ctx *parser.GrantRoleToRoleContext) {}

func (l *Listener) EnterRevokeRoleFromUser(ctx *parser.RevokeRoleFromUserContext) {}

func (l *Listener) EnterRevokeRoleFromRole(ctx *parser.RevokeRoleFromRoleContext) {}
