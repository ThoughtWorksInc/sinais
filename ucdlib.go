package runefinder

import ( // <1>
	"strconv"
	"strings"
)

func AnalisarLinha(ucdLine string) (rune, string) {
	campos := strings.Split(ucdLine, ";")            // <2>
	código, _ := strconv.ParseInt(campos[0], 16, 32) // <3>
	return rune(código), campos[1]                   // <4>
}
