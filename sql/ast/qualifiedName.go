// Package ast @author: Violet-Eva @date  : 2025/5/7 @notes :
package ast

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/violet-eva-01/StarRocks/sql/parser"
	"strings"
)

type QualifiedName struct {
	catalog  string
	database string
	parts    []string
}

// NewQualifiedName 创建 QualifiedName 实例
func NewQualifiedName(originalParts []string) QualifiedName {
	if originalParts == nil {
		panic("originalParts is null")
	}
	if len(originalParts) == 0 {
		panic("originalParts is empty")
	}
	// 复制切片以确保不可变性
	partsCopy := make([]string, len(originalParts))
	copy(partsCopy, originalParts)
	return QualifiedName{
		parts: partsCopy,
	}
}

func GetQualifiedName(ctx antlr.ParserRuleContext) QualifiedName {
	var (
		parts []string
	)
	for _, v := range ctx.GetChildren() {
		switch v.(type) {
		case antlr.TerminalNode:
			node := v.(antlr.TerminalNode)
			if node.GetSymbol().GetTokenType() == parser.StarRocksParserDOT_IDENTIFIER {
				parts = append(parts, node.GetText()[1:])
			}
		case parser.IIdentifierOrStringOrStarContext:
			identifier := v.(parser.IIdentifierOrStringOrStarContext).Identifier()
			parts = append(parts, identifier.GetText())
		case parser.IIdentifierContext:
			identifier := v.(parser.IIdentifierContext)
			parts = append(parts, identifier.GetText())
		}
	}

	return NewQualifiedName(parts)
}

func (qn QualifiedName) GetParts() []string {
	return qn.parts
}

func (qn QualifiedName) String() string {
	return strings.Join(qn.parts, ".")
}

func (qn QualifiedName) ToSql() string {
	return qn.String()
}
