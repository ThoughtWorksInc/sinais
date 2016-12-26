package runefinder

import "testing"

const linhaLetraA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

func TestAnalisarLinha(t *testing.T) {
	runa, nome := AnalisarLinha(linhaLetraA)
	if runa != 'A' {
		t.Errorf("Esperava 'A', veio %q", runa) // <1>
	}
	const nomeA = "LATIN CAPITAL LETTER A" // <2>
	if nome != nomeA {
		t.Errorf("Esperava %q, veio %q", nomeA, nome) // <3>
	}
}
