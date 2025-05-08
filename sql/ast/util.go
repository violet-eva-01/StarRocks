// Package ast @author: Violet-Eva @date  : 2025/5/8 @notes :
package ast

func findIndex(str string, strArr []string) int {
	for index, element := range strArr {
		if str == element {
			return index
		}
	}
	return -1
}