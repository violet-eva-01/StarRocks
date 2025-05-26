// Package parser @author: Violet-Eva @date  : 2025/5/25 @notes :
package sql

import (
	"testing"
)

func TestRevoke(t *testing.T) {
	query := `
set @policyName="Revoke-20250521-0001";
revoke all PRIVILEGES,select,alter,create materialized view on alden.test,alden2.test2 from 'violet'@'192.168.%';
`
	_ = NewParser(query, "", "")
}

func TestGrant(t *testing.T) {
	query := "set @policyName=\"Revoke-20250521-0001\";grant all PRIVILEGES,select,alter,create materialized view on alden.test,alden2.test2 to 'violet'@'192.168.%' with grant option;grant all on alden to role aldentest;grant select on alden.test to violet;"
	_ = NewParser(query, "", "")
}

func TestExecuteAs(t *testing.T) {
	query := "execute as violet with no revert"
	_ = NewParser(query, "", "")
}