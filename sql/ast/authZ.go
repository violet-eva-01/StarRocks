// Package ast @author: Violet-Eva @date  : 2025/5/21 @notes :
package ast

import (
	"github.com/violet-eva-01/StarRocks/sql/parser"
)
type AuthZType int

const (
	AZGrant AuthZType = iota
	AZRevoke
)

var authZTypeName = map[string]AuthZType{
	"GRANT":  AZGrant,
	"REVOKE": AZRevoke,
}

func (t AuthZType) String() string {
	return findName(t, authZTypeName)
}

type AuthZ struct {
	Catalog      string
	Type         AuthZType
	Permissions  map[PrivilegeType]CreatePrivilegeType
	Object       Object
	UserIdentity UserIdentity
	Option       bool
}

func NewEmptyAuthZ() *AuthZ {
	return &AuthZ{}
}

func NewAuthZWithGrantOnTableBrief(catalog, database string, ctx *parser.GrantOnTableBriefContext) []AuthZ {
	var authZs []AuthZ
	for _, i := range GetObjectNameList(catalog, database, ctx.PrivObjectNameList()) {
		authZ := NewEmptyAuthZ()
		authZ.Catalog = catalog
		authZ.Type = AZGrant
		authZ.Permissions = GetTypeList(ctx.PrivilegeTypeList())
		authZ.Object = *GetObject(i, "")
		authZ.UserIdentity = *GetUserIdentify(ctx.GrantRevokeClause())
		if ctx.WITH() != nil && ctx.OPTION() != nil && ctx.GRANT(1) != nil {
			authZ.Option = true
		}
		authZs = append(authZs, *authZ)
	}
	return nil
}

func NewAuthZWithRevokeOnTableBrief(catalog, database string, ctx *parser.RevokeOnTableBriefContext) *AuthZ {
	/*	ptList := GetTypeList(ctx.PrivilegeTypeList())
		tns := GetObjectNameList(catalog, database, ctx.PrivObjectNameList())
		userIdentify := GetUserIdentify(ctx.GrantRevokeClause())
		return &AuthZ{
			Catalog:          catalog,
			Type:             AZRevoke,
			Permissions:      ptList,
			ObjectType:       -1,
			ObjectTypePlural: -1,
			TableName:        tns,
			UserIdentity:     *userIdentify,
		}*/
	return nil
}
