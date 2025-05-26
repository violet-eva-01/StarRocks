// Package ast @author: Violet-Eva @date  : 2025/5/8 @notes :
package ast

type t interface {
	PrivilegeType | CreatePrivilegeType | PrivilegeObjectType | PrivilegeObjectTypePlural | AuthZType | privilegeOnType
}

func findName[T t](input T, slice map[string]T) string {
	for k, v := range slice {
		if v == input {
			return k
		}
	}
	return ""
}
