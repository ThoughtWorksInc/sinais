package runefinder // <1>

import "testing" // <2>

const linhaLetraA = `0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;` // <3>

func TestAnalisarLinha(t *testing.T) { // <4>
	runa, _ := AnalisarLinha(linhaLetraA) // <5>
	if runa != 'A' {                      // <6>
		t.Errorf("Esperava 'A', veio %c", runa) // <7>
	}
}
