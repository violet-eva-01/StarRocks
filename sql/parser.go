// Package sql @author: Violet-Eva @date  : 2025/5/8 @notes :
package sql

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/violet-eva-01/StarRocks/sql/ast"
	"github.com/violet-eva-01/StarRocks/sql/parser"
)

func NewParser(query string, catalog string, dbName string) (al *ast.Listener) {
	listener := ast.NewListener()
	starRocksStream := antlr.NewInputStream(query)
	lexer := parser.NewStarRocksLexer(starRocksStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	starRocksParser := parser.NewStarRocksParser(tokenStream)
	if catalog != "" {
		listener.SetCatalog(catalog)
	}
	if dbName != "" {
		listener.SetDatabase(dbName)
	}
	starRocksParser.RemoveErrorListeners()
	starRocksParser.AddErrorListener(listener)
	defer func() {
		if err := recover(); err != nil {
			al = listener
		}
	}()
	antlr.ParseTreeWalkerDefault.Walk(listener, starRocksParser.SqlStatements())
	return listener
}
