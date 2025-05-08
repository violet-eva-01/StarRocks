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
	pos      NodePosition
}

// NewQualifiedName 创建 QualifiedName 实例，不带位置参数
func NewQualifiedName(originalParts []string) QualifiedName {
	return NewQualifiedNameWithPos(originalParts, Zero)
}

// NewQualifiedNameWithPos 创建 QualifiedName 实例，带位置参数
func NewQualifiedNameWithPos(originalParts []string, pos NodePosition) QualifiedName {
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
		pos:   pos,
	}
}

func GetQualifiedName(ctx antlr.ParserRuleContext) QualifiedName {
	var (
		parts []string
	)
	pos := CreatePos(ctx)
	for _, v := range ctx.GetChildren() {
		switch v.(type) {
		case antlr.TerminalNode:
			node := v.(antlr.TerminalNode)
			if node.GetSymbol().GetTokenType() == parser.StarRocksParserDOT_IDENTIFIER {
				parts = append(parts, node.GetText()[1:])
			}
		case parser.IIdentifierContext:
			context := v.(parser.IIdentifierContext)
			identifier := context
			parts = append(parts, identifier.GetText())
		}
	}

	return NewQualifiedNameWithPos(parts, pos)
}

// GetParts 获取 parts 字段
func (qn QualifiedName) GetParts() []string {
	return qn.parts
}

// String 实现 Stringer 接口
func (qn QualifiedName) String() string {
	return strings.Join(qn.parts, ".")
}

// GetPos 获取 pos 字段
func (qn QualifiedName) GetPos() NodePosition {
	return qn.pos
}

// ToSql 实现 ParseNode 接口的 ToSql 方法
func (qn QualifiedName) ToSql() string {
	return qn.String()
}
