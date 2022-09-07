package main

import "regexp"

var (
	removeReg = make([]*regexp.Regexp, 0)
	moveReg   = make([]*regexp.Regexp, 0)
)

func InitRegex(remove []string, move []string) {
	for i := range remove {
		removeReg = append(removeReg, regexp.MustCompile(remove[i]))
	}
	for i := range move {
		moveReg = append(moveReg, regexp.MustCompile(move[i]))
	}
}
