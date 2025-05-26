// Package ast @author: Violet-Eva @date  : 2025/5/22 @notes :
package ast

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/violet-eva-01/StarRocks/sql/parser"
	"strings"
)

const (
	defaultHost = "%"
	defaultForm = "USER"
)

type UserIdentity struct {
	form string
	name string
	host string
}

func NewUserIdentity(form, name, host string) *UserIdentity {
	if form == "" {
		form = defaultForm
	}
	return &UserIdentity{form, name, host}
}

func GetUserIdentify(ctx antlr.ParserRuleContext) *UserIdentity {
	var (
		name string
		host string
		form string
	)
	for _, i := range ctx.GetChildren() {
		switch i.(type) {
		case *parser.UserWithHostContext:
			user := i.(*parser.UserWithHostContext)
			name = strings.Trim(user.IdentifierOrString(0).GetText(), "\"'")
			host = strings.Trim(user.IdentifierOrString(1).GetText(), "\"'")
		case *parser.UserWithoutHostContext:
			user := i.(*parser.UserWithoutHostContext)
			name = strings.Trim(user.IdentifierOrString().GetText(), "\"'")
			host = defaultHost
		case *parser.IdentifierOrStringContext:
			name = i.(*parser.IdentifierOrStringContext).Identifier().GetText()
		case antlr.TerminalNode:
			form = strings.ToUpper(i.(antlr.TerminalNode).GetSymbol().GetText())
		}
	}
	return NewUserIdentity(form, name, host)
}

func (u *UserIdentity) SetHost(host string) {
	u.host = host
}

func (u *UserIdentity) SetForm(form string) {
	u.form = form
}

func (u *UserIdentity) SetName(name string) {
	u.name = name
}

func (u *UserIdentity) GetHost() string {
	return u.host
}

func (u *UserIdentity) GetForm() string {
	return u.form
}

func (u *UserIdentity) GetName() string {
	return u.name
}

func (u *UserIdentity) String() string {
	return u.name
}

func (u *UserIdentity) ToSql() string {
	var sb strings.Builder
	sb.WriteString("'")
	sb.WriteString(u.name)
	sb.WriteString("'")
	if u.host != "" {
		sb.WriteString("@")
		sb.WriteString("'")
		sb.WriteString(u.host)
		sb.WriteString("'")
	}
	return sb.String()
}
