package ast

import (
	"strings"
)

type Action int

const (
	CREATE Action = 0 + iota
	ALTER
	TRUNCATE
	DROP
	SELECT
	INSERT
	UPDATE
	DELETE
)

var actionNames = []string{
	"SELECT",
	"ALTER",
	"INSERT",
	"CREATE",
	"DROP",
	"UPDATE",
	"DELETE",
	"TRUNCATE",
}

func GetActionFromString(str string) Action {

	index := findIndex(strings.ToUpper(str), actionNames)
	if index == -1 {
		return -1
	} else {
		return Action(index)
	}

}

func (a Action) String() string {

	if a >= CREATE && a <= DELETE {
		return actionNames[a]
	}

	return "nil"
}
