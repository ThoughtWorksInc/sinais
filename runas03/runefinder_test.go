package main

import (
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
	runa, nome := AnalisarLinha(linhaLetraA)
	if runa != 'A' {
		t.Errorf("Esperava 'A', veio %q", runa)
	}
	const nomeA = "LATIN CAPITAL LETTER A"
	if nome != nomeA {
		t.Errorf("Esperava %q, veio %q", nomeA, nome)
	}
}

func ExampleListar() { // ➊
	texto := strings.NewReader(linhas3Da43) // ➋
	Listar(texto, "MARK")                   // ➌
	// Output: U+003F	?	QUESTION MARK
}

func ExampleListar_doisResultados() { // ➊
	texto := strings.NewReader(linhas3Da43)
	Listar(texto, "SIGN") // ➋
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}
