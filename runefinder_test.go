package main

import (
	"os"
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

func TestListar(t *testing.T) {
	oldStdOut := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() {
		os.Stdout.Close()
		os.Stdout = oldStdOut
	}()

	texto := strings.NewReader(linhas3Da43)
	count := Listar(texto, "MARK")
	expectedCount := 1
	if count != expectedCount {
		t.Errorf("Esperava %d, veio %d", expectedCount, count)
	}
}

func ExampleListar() {
	texto := strings.NewReader(linhas3Da43)
	Listar(texto, "MARK")
	// Output: U+003F	?	QUESTION MARK
}

func ExampleListar_doisResultados() { // <1>
	texto := strings.NewReader(linhas3Da43)
	Listar(texto, "SIGN") // <2>
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}

func Example() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "cruzeiro"}
	main()
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
	// 1 character found
}

func Example_tresResultados() {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "pound"}
	main()
	// Output:
	// U+00A3	£	POUND SIGN
	// U+FFE1	￡	FULLWIDTH POUND SIGN
	// U+1F4B7	💷	BANKNOTE WITH POUND SIGN
	// 3 characters found
}
