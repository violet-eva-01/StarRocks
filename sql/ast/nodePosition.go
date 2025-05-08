// Package ast @author: Violet-Eva @date  : 2025/5/7 @notes :
package ast

import "github.com/antlr4-go/antlr/v4"

// NodePosition 用于记录 SQL 中元素的位置
type NodePosition struct {
	Line    int
	Col     int
	EndLine int
	EndCol  int
}

// Zero 表示零位置
var Zero = NodePosition{0, 0, 0, 0}

func CreatePos(ctx antlr.ParserRuleContext) NodePosition {
	start := ctx.GetStart()
	stop := ctx.GetStop()
	if start == nil {
		return Zero
	}
	if stop == nil {
		return NewNodePosition(start.GetLine(), start.GetColumn())
	}
	return NewNodePositionFromTokens(start, stop)
}

// NewNodePositionFromTerminalNode 根据 TerminalNode 创建 NodePosition
func NewNodePositionFromTerminalNode(node antlr.TerminalNode) NodePosition {
	symbol := node.GetSymbol()
	return NewNodePosition(symbol.GetLine(), symbol.GetColumn())
}

// NewNodePositionFromToken 根据 Token 创建 NodePosition
func NewNodePositionFromToken(token antlr.Token) NodePosition {
	return NewNodePosition(token.GetLine(), token.GetLine())
}

// NewNodePositionFromTokens 根据起始和结束 Token 创建 NodePosition
func NewNodePositionFromTokens(start, end antlr.Token) NodePosition {
	return NodePosition{
		Line:    start.GetLine(),
		Col:     start.GetColumn(),
		EndLine: end.GetLine(),
		EndCol:  end.GetColumn(),
	}
}

// NewNodePosition 创建一个新的 NodePosition
func NewNodePosition(line, col int) NodePosition {
	return NodePosition{
		Line:    line,
		Col:     col,
		EndLine: line,
		EndCol:  col,
	}
}

// GetLine 获取起始行号
func (np NodePosition) GetLine() int {
	return np.Line
}

// GetCol 获取起始列号
func (np NodePosition) GetCol() int {
	return np.Col
}

// GetEndLine 获取结束行号
func (np NodePosition) GetEndLine() int {
	return np.EndLine
}

// GetEndCol 获取结束列号
func (np NodePosition) GetEndCol() int {
	return np.EndCol
}

// IsZero 判断是否为零位置
func (np NodePosition) IsZero() bool {
	return np.Line == 0 && np.Col == 0 && np.EndLine == 0 && np.EndCol == 0
}
