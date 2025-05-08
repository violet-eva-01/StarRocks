// Package sql @author: Violet-Eva @date  : 2025/5/8 @notes :
package sql

import (
	"fmt"
	"testing"
)

func TestNewParser(t *testing.T) {
	query := "create table ads.test as select * from ads.test1 left join ads.test2 on ads.test1.id=ads.test2.id;"
	l := NewParser(query, "dc", "ddb")
	names := l.GetActionTables()
	for action, name := range names {
		for _, tbl := range name {
			fmt.Println(action, tbl.String())
		}
	}

}
