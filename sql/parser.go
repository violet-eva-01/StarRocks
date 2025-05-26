// Package sql @author: Violet-Eva @date  : 2025/5/8 @notes :
package sql

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/violet-eva-01/StarRocks/sql/ast"
	"github.com/violet-eva-01/StarRocks/sql/parser"
)

func NewParser(query string, catalog string, dbName string) *ast.Listener {
	starRocksStream := antlr.NewInputStream(query)
	lexer := parser.NewStarRocksLexer(starRocksStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	starRocksParser := parser.NewStarRocksParser(tokenStream)
	listener := ast.NewListener()
	if catalog != "" {
		listener.SetCatalog(catalog)
	}
	if dbName != "" {
		listener.SetDatabase(dbName)
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, starRocksParser.SqlStatements())
	return listener
}
