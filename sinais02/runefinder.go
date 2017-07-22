package main

import ( // ➊
	"strconv"
	"strings"
)

func AnalisarLinha(ucdLine string) (rune, string) {
	campos := strings.Split(ucdLine, ";")            // ➋
	código, _ := strconv.ParseInt(campos[0], 16, 32) // ➌
	return rune(código), campos[1]                   // ➍
}
