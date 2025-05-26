// Package ast @author: Violet-Eva @date  : 2025/5/8 @notes :
package ast

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/violet-eva-01/StarRocks/sql/parser"
	"strings"
)

var (
	//  isCaseSensitive -> TableName
	isCaseSensitive = false
)

func SetCaseSensitive() {
	isCaseSensitive = true
}

const (
	defaultCatalog  = "default_catalog"
	defaultDatabase = "default_db"
	defaultGrantVar = "policyName"
)

type ListenerType int

const (
	listenerSql ListenerType = iota
	listenerAuthZ
	listenerAudit
)

type Listener struct {
	Type            ListenerType
	defaultCatalog  string
	defaultDatabase string
	withTable       []TableName
	UserProperties  map[string]string
	SystemVariables map[string]string
	Tables          []TableName
	ActionTables    map[PrivilegeType]map[TableName]int
	defaultGrantVar string
	policyName      string
	authZs          []AuthZ
	*parser.BaseStarRocksListener
}

func NewListener() *Listener {
	return &Listener{
		Type:            listenerSql,
		defaultCatalog:  defaultCatalog,
		defaultDatabase: defaultDatabase,
		defaultGrantVar: defaultGrantVar,
		ActionTables:    make(map[PrivilegeType]map[TableName]int),
		UserProperties:  make(map[string]string),
		SystemVariables: make(map[string]string),
	}
}

func (l *Listener) SetCatalog(catalog string) {
	l.defaultCatalog = catalog
}

func (l *Listener) SetDatabase(database string) {
	l.defaultDatabase = database
}

func (l *Listener) SetGrantVar(grantVar string) {
	l.defaultGrantVar = grantVar
}

func (l *Listener) GetCatalog() string {
	return l.defaultCatalog
}

func (l *Listener) GetDatabase() string {
	return l.defaultDatabase
}

func (l *Listener) GetGrantVar() string { return l.defaultGrantVar }

func (l *Listener) GetTables() []TableName {
	return l.Tables
}

func (l *Listener) GetActionTables() map[PrivilegeType]map[TableName]int {
	return l.ActionTables
}

func (l *Listener) GetWithTables() []TableName {
	return l.withTable
}

func (l *Listener) GetPolicyName() string { return l.policyName }

func (l *Listener) GetTableNamesByAction(action PrivilegeType) map[TableName]int {
	return l.ActionTables[action]
}

func (l *Listener) GetUserProperties() map[string]string {
	return l.UserProperties
}

func (l *Listener) GetSystemVariables() map[string]string {
	return l.SystemVariables
}

func (l *Listener) GetUserPropertyByName(str string) string {
	if l.UserProperties != nil && l.UserProperties[str] != "" {
		return l.UserProperties[str]
	}
	return ""
}

func (l *Listener) GetSystemVariableByName(str string) string {
	if l.SystemVariables != nil && l.SystemVariables[str] != "" {
		return l.SystemVariables[str]
	}
	return ""
}

func (l *Listener) GetTableNameWithAction(action string, ctx antlr.ParserRuleContext) {
	l.GetTableNameWithQualifiedName(action, GetQualifiedName(ctx))
}

func (l *Listener) GetTableNameWithQualifiedName(action string, qf QualifiedName) {
	tbl := NewTableNameWithQualifiedName(l.defaultCatalog, l.defaultDatabase, qf)
	l.GetTableNameWithTableName(action, *tbl)
}

func (l *Listener) GetTableNameWithTableName(strAction string, tbl TableName) {
	action := ParsePrivilegeType(strAction)
	if !tbl.InTableNameList(l.withTable) {
		if l.ActionTables[action] == nil && l.ActionTables[action][tbl] == 0 {
			ti := make(map[TableName]int)
			ti[tbl] = 1
			l.ActionTables[action] = ti
		} else {
			l.ActionTables[action][tbl] += 1
		}
		l.Tables = append(l.Tables, tbl)
	}
}

// Analyze Statement

func (l *Listener) EnterAnalyzeStatement(ctx *parser.AnalyzeStatementContext) {
	// analyze 归属于 insert 权限
	if ctx.TableName() != nil {
		l.GetTableNameWithAction("INSERT", ctx.TableName().QualifiedName())
	}
}

func (l *Listener) EnterCreateAnalyzeStatement(ctx *parser.CreateAnalyzeStatementContext) {
	if ctx.TABLE() != nil {
		l.GetTableNameWithAction("INSERT", ctx.QualifiedName(0))
	}
}

func (l *Listener) EnterDropStatsStatement(ctx *parser.DropStatsStatementContext) {
	l.GetTableNameWithAction("INSERT", ctx.QualifiedName())
}

func (l *Listener) EnterHistogramStatement(ctx *parser.HistogramStatementContext) {
	l.GetTableNameWithAction("INSERT", ctx.TableName().QualifiedName())
}

func (l *Listener) EnterDropHistogramStatement(ctx *parser.DropHistogramStatementContext) {
	if ctx.TABLE() != nil {
		l.GetTableNameWithAction("INSERT", ctx.QualifiedName(0))
	}
}

// Catalog Statement

/*
   | createExternalCatalogStatement
   | dropExternalCatalogStatement
   | showCatalogsStatement
   | showCreateExternalCatalogStatement
   | alterCatalogStatement
*/

func (l *Listener) EnterSetCatalogStatement(ctx *parser.SetCatalogStatementContext) {
	catalog := strings.ToLower(ctx.IdentifierOrString().Identifier().GetText())
	l.defaultCatalog = catalog
}

// Database Statement

func (l *Listener) EnterUseDatabaseStatement(ctx *parser.UseDatabaseStatementContext) {
	parts := GetQualifiedName(ctx.QualifiedName()).GetParts()
	if len(parts) == 2 {
		l.defaultCatalog = strings.ToLower(parts[0])
		l.defaultDatabase = strings.ToLower(parts[1])
	} else if len(parts) == 1 {
		l.defaultDatabase = strings.ToLower(parts[0])
	}
}

/*
   | alterDbQuotaStatement
   | createDbStatement
   | dropDbStatement
   | alterDatabaseRenameStatement
   | recoverDbStmt
*/

// Function Statement

/*
   | dropFunctionStatement
   | createFunctionStatement
*/

// Routine Statement

func (l *Listener) EnterCreateRoutineLoadStatement(ctx *parser.CreateRoutineLoadStatementContext) {
	l.GetTableNameWithAction("INSERT", ctx.GetTable())
}

// Load Statement

func (l *Listener) EnterLoadStatement(ctx *parser.LoadStatementContext) {
	db := ctx.GetLabel().GetDb().GetText()
	for _, dataDesc := range ctx.GetData().AllDataDesc() {
		if src := dataDesc.GetSrcTableName(); src != nil {
			// hive外表(properties 指定hive表属性)需要创建在内表同库下
			l.GetTableNameWithTableName("SELECT", *NewTableName(l.defaultCatalog, db, src.GetText()))
		}
		if dst := dataDesc.GetDstTableName(); dst != nil {
			l.GetTableNameWithTableName("INSERT", *NewTableName(l.defaultCatalog, db, dst.GetText()))
		}
	}
}

// Table DDL Statement

/*
   | exportStatement
*/

func (l *Listener) EnterExportStatement(ctx *parser.ExportStatementContext) {
	l.GetTableNameWithAction(ctx.EXPORT().GetText(), ctx.TableDesc().QualifiedName())
}

func (l *Listener) EnterCreateTableStatement(ctx *parser.CreateTableStatementContext) {
	l.GetTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterCreateTableAsSelectStatement(ctx *parser.CreateTableAsSelectStatementContext) {
	l.GetTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())

}

func (l *Listener) EnterCreateTableLikeStatement(ctx *parser.CreateTableLikeStatementContext) {
	l.GetTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName(0))
	l.GetTableNameWithAction("SELECT", ctx.QualifiedName(1))
}

func (l *Listener) EnterShowCreateTableStatement(ctx *parser.ShowCreateTableStatementContext) {
	// show 不需要权限暂定为select一次表
	l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
}

func (l *Listener) EnterAlterTableStatement(ctx *parser.AlterTableStatementContext) {
	l.GetTableNameWithAction(ctx.ALTER().GetText(), ctx.QualifiedName())
	for _, alterClause := range ctx.AllAlterClause() {
		if alterClause.TableRenameClause() != nil {
			l.GetTableNameWithQualifiedName("CREATE", GetQualifiedName(alterClause.TableRenameClause()))
		}
	}
}

func (l *Listener) EnterShowAlterStatement(ctx *parser.ShowAlterStatementContext) {
	if ctx.FROM() != nil {
		l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
	}
}

func (l *Listener) EnterCancelAlterTableStatement(ctx *parser.CancelAlterTableStatementContext) {
	l.GetTableNameWithAction(ctx.ALTER().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterCreateIndexStatement(ctx *parser.CreateIndexStatementContext) {
	// create index 归属于 alter权限
	l.GetTableNameWithAction("ALTER", ctx.QualifiedName())
}

func (l *Listener) EnterDropIndexStatement(ctx *parser.DropIndexStatementContext) {
	// drop index 归属于 alter权限
	l.GetTableNameWithAction("ALTER", ctx.QualifiedName())
}

func (l *Listener) EnterShowIndexStatement(ctx *parser.ShowIndexStatementContext) {
	tbl := NewTableNameWithQualifiedName(l.defaultCatalog, l.defaultDatabase, GetQualifiedName(ctx.QualifiedName(0)))
	if ctx.FROM(1) != nil {
		tbl.SetDb(ctx.QualifiedName(1).GetText())
	}
	l.GetTableNameWithTableName("SELECT", *tbl)
}

func (l *Listener) EnterTruncateTableStatement(ctx *parser.TruncateTableStatementContext) {
	// truncate 归属 delete 权限
	l.GetTableNameWithAction("DELETE", ctx.QualifiedName())
}

func (l *Listener) EnterDropTableStatement(ctx *parser.DropTableStatementContext) {
	l.GetTableNameWithAction(ctx.DROP().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterRecoverTableStatement(ctx *parser.RecoverTableStatementContext) {
	// recover table 归属于 create 权限
	l.GetTableNameWithAction("CREATE", ctx.QualifiedName())
}

func (l *Listener) EnterRecoverPartitionStatement(ctx *parser.RecoverPartitionStatementContext) {
	// recover partition 归属于 insert 权限
	l.GetTableNameWithAction("INSERT", ctx.QualifiedName())
}

func (l *Listener) EnterShowPartitionsStatement(ctx *parser.ShowPartitionsStatementContext) {
	if ctx.FROM() != nil {
		l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
	}
}

func (l *Listener) EnterShowDynamicPartitionStatement(ctx *parser.ShowDynamicPartitionStatementContext) {
	if ctx.FROM() != nil {
		l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
	}
}

func (l *Listener) EnterShowColumnStatement(ctx *parser.ShowColumnStatementContext) {
	if ctx.GetTable() != nil {
		l.GetTableNameWithAction("SELECT", ctx.GetTable())
	}
}

func (l *Listener) EnterShowTabletStatement(ctx *parser.ShowTabletStatementContext) {
	if ctx.FROM() != nil {
		l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
	}
}

func (l *Listener) EnterRefreshTableStatement(ctx *parser.RefreshTableStatementContext) {
	l.GetTableNameWithAction(ctx.REFRESH().GetText(), ctx.QualifiedName())
}

// View Statement

func (l *Listener) EnterCreateViewStatement(ctx *parser.CreateViewStatementContext) {
	l.GetTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterDropViewStatement(ctx *parser.DropViewStatementContext) {
	l.GetTableNameWithAction(ctx.DROP().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterAlterViewStatement(ctx *parser.AlterViewStatementContext) {
	l.GetTableNameWithAction(ctx.ALTER().GetText(), ctx.QualifiedName())
}

// Materialized View Statement

func (l *Listener) EnterCreateMaterializedViewStatement(ctx *parser.CreateMaterializedViewStatementContext) {
	l.GetTableNameWithAction(ctx.CREATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterDropMaterializedViewStatement(ctx *parser.DropMaterializedViewStatementContext) {
	l.GetTableNameWithAction(ctx.DROP().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterAlterMaterializedViewStatement(ctx *parser.AlterMaterializedViewStatementContext) {
	l.GetTableNameWithAction(ctx.ALTER().GetText(), ctx.QualifiedName())
	if ctx.TableRenameClause() != nil {
		l.GetTableNameWithQualifiedName("CREATE", GetQualifiedName(ctx.TableRenameClause()))
	}
}

func (l *Listener) EnterRefreshMaterializedViewStatement(ctx *parser.RefreshMaterializedViewStatementContext) {
	l.GetTableNameWithAction(ctx.REFRESH().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterCancelRefreshMaterializedViewStatement(ctx *parser.CancelRefreshMaterializedViewStatementContext) {
	l.GetTableNameWithAction(ctx.REFRESH().GetText(), ctx.QualifiedName())
}

// DML Statement

func (l *Listener) EnterWithClause(ctx *parser.WithClauseContext) {
	for _, v := range ctx.AllCommonTableExpression() {
		withTBL := NewTableNameWithQualifiedName(l.defaultCatalog, l.defaultDatabase, GetQualifiedName(v))
		l.withTable = append(l.withTable, *withTBL)
	}
}

func (l *Listener) EnterDescTableStatement(ctx *parser.DescTableStatementContext) {
	l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
}

func (l *Listener) EnterTableAtom(ctx *parser.TableAtomContext) {
	l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
}

func (l *Listener) EnterInsertStatement(ctx *parser.InsertStatementContext) {
	l.GetTableNameWithAction(ctx.INSERT().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterUpdateStatement(ctx *parser.UpdateStatementContext) {
	l.GetTableNameWithAction(ctx.UPDATE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterDeleteStatement(ctx *parser.DeleteStatementContext) {
	l.GetTableNameWithAction(ctx.DELETE().GetText(), ctx.QualifiedName())
}

func (l *Listener) EnterShowDeleteStatement(ctx *parser.ShowDeleteStatementContext) {
	if ctx.FROM() != nil {
		l.GetTableNameWithAction("SELECT", ctx.QualifiedName())
	}
}

// authZs Statement

func GetAuthZ(policyName, catalog, database string, ctx antlr.ParserRuleContext) []AuthZ {
	var (
		authZ []AuthZ
	)
	// [role|role slice]
	// [user|user slice]
	// [table|table slice|all tables]
	// [function|function slice|all functions]
	// [system]
	// []
	switch ctx.(type) {
	case *parser.GrantRoleToUserContext:
	case *parser.RevokeRoleFromUserContext:
	case *parser.GrantRoleToRoleContext:
	case *parser.RevokeRoleFromRoleContext:
	case *parser.GrantOnUserContext:
	case *parser.RevokeOnUserContext:
	case *parser.GrantOnTableBriefContext:
	case *parser.RevokeOnTableBriefContext:
	case *parser.GrantOnFuncContext:
	case *parser.RevokeOnFuncContext:
	case *parser.GrantOnSystemContext:
	case *parser.RevokeOnSystemContext:
	case *parser.GrantOnPrimaryObjContext:
	case *parser.RevokeOnPrimaryObjContext:
	case *parser.GrantOnAllContext:
	case *parser.RevokeOnAllContext:
	}
	return authZ
}

func (l *Listener) EnterCreateUserStatement(ctx *parser.CreateUserStatementContext) {
}

func (l *Listener) EnterDropUserStatement(ctx *parser.DropUserStatementContext) {
}

func (l *Listener) EnterAlterUserStatement(ctx *parser.AlterUserStatementContext) {

}

// set @policyName="Grant-xxx";GRANT xxx;
// set @policyName="Revoke-xxx";REVOKE xxx;

func (l *Listener) EnterSetUserVar(ctx *parser.SetUserVarContext) {
	variable := ctx.UserVariable().IdentifierOrString().GetText()
	value := strings.Trim(ctx.Expression().GetChild(0).(*parser.BooleanExpressionDefaultContext).
		Predicate().Get_valueExpression().GetText(), "\"'")
	if variable == l.defaultGrantVar {
		l.policyName = value
	}
	l.UserProperties[variable] = value
}

func (l *Listener) EnterSetSystemVar(ctx *parser.SetSystemVarContext) {
	variable := ctx.Identifier().GetText()
	value := ctx.SetExprOrDefault().Expression().GetText()
	l.SystemVariables[variable] = value
}

func (l *Listener) EnterGrantRoleToUser(ctx *parser.GrantRoleToUserContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterGrantRoleToRole(ctx *parser.GrantRoleToRoleContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeRoleFromUser(ctx *parser.RevokeRoleFromUserContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeRoleFromRole(ctx *parser.RevokeRoleFromRoleContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterGrantOnUser(ctx *parser.GrantOnUserContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterGrantOnTableBrief(ctx *parser.GrantOnTableBriefContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterGrantOnFunc(ctx *parser.GrantOnFuncContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterGrantOnSystem(ctx *parser.GrantOnSystemContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterGrantOnPrimaryObj(ctx *parser.GrantOnPrimaryObjContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterGrantOnAll(ctx *parser.GrantOnAllContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeOnUser(ctx *parser.RevokeOnUserContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeOnTableBrief(ctx *parser.RevokeOnTableBriefContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeOnFunc(ctx *parser.RevokeOnFuncContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeOnSystem(ctx *parser.RevokeOnSystemContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeOnPrimaryObj(ctx *parser.RevokeOnPrimaryObjContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

func (l *Listener) EnterRevokeOnAll(ctx *parser.RevokeOnAllContext) {
	l.Type = listenerAuthZ
	l.authZs = append(l.authZs, GetAuthZ(l.policyName, l.defaultCatalog, l.defaultDatabase, ctx)...)
}

// High Risk

func (l *Listener) EnterExecuteAsStatement(ctx *parser.ExecuteAsStatementContext) {
	/*	_ = ExecuteAs
		user := ctx.User().GetText()*/
}
