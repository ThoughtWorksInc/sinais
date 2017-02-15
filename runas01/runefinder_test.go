package main // ➊

import "testing" // ➋

const linhaLetraA = `0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;` // ➌

func TestAnalisarLinha(t *testing.T) { // ➍
	runa, _ := AnalisarLinha(linhaLetraA) // ➎
	if runa != 'A' {                      // ➏
		t.Errorf("Esperava 'A', veio %c", runa) // ➐
	}
}
