package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

const linhaLetraA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

const linhas3Da43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestAnalisarLinha(t *testing.T) {
	var casos = []struct { // ➊
		linha    string
		runa     rune
		nome     string
		palavras []string
	}{ // ➋
		{"0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;",
			'!', "EXCLAMATION MARK", []string{"EXCLAMATION", "MARK"}},
		{"002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;",
			'-', "HYPHEN-MINUS", []string{"HYPHEN", "MINUS"}},
		{"0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;",
			'\'', "APOSTROPHE (APOSTROPHE-QUOTE)", []string{"APOSTROPHE", "QUOTE"}},
	}
	for _, caso := range casos { // ➌
		runa, nome, palavras := AnalisarLinha(caso.linha) // ➍
		if runa != caso.runa || nome != caso.nome ||
			!reflect.DeepEqual(palavras, caso.palavras) {
			t.Errorf("\nAnalisarLinha(%q)\n-> (%q, %q, %q)", // ➎
				caso.linha, runa, nome, palavras)
		}
	}
}

func ExampleListar() {
	texto := strings.NewReader(linhas3Da43)
	Listar(texto, "MARK")
	// Output: U+003F	?	QUESTION MARK
}

func ExampleListar_doisResultados() {
	texto := strings.NewReader(linhas3Da43)
	Listar(texto, "SIGN")
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}

func Example() { // ➊
	oldArgs := os.Args                   // ➋
	defer func() { os.Args = oldArgs }() // ➌
	os.Args = []string{"", "cruzeiro"}   // ➍
	main()                               // ➎
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
}
