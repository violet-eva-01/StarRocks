// Package sql @author: Violet-Eva @date  : 2025/5/8 @notes :
package sql

import (
	"fmt"
	"testing"
	"time"
)

func TestNewParser(t *testing.T) {
	query := "create table ads.test as select * from ads.test1 left join ads.test2 on ads.test1.id=ads.test2.id 11;"
	l := NewParser(query, "dc", "ddb")
	names := l.GetActionTables()
	for action, name := range names {
		for _, tbl := range name {
			fmt.Println(action, tbl.String())
		}
	}

}

func TestName(t *testing.T) {
	ticker := time.NewTicker(15 * time.Second)
	timer := time.NewTimer(0)
	defer ticker.Stop()
	fmt.Println(time.Now().Format(time.DateTime))
	for {
		select {
		case <-timer.C:
			fmt.Println(time.Now().Format(time.DateTime), "timer")
			timer.Stop()
		case <-ticker.C:
			fmt.Println(time.Now().Format(time.DateTime), "tick")
		}
	}
}
