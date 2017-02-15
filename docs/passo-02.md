---
permalink: passo-02
---

# Runas, passo 2: primeira função completa

Agora vamos terminar de implementar a função `AnalisarLinha`, guiados por testes.

Alteramos `runefinder_test.go` para verificar o nome devolvido por `AnalisarLinha`:

```go
package main

import "testing"

const linhaLetraA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

func TestAnalisarLinha(t *testing.T) {
	runa, nome := AnalisarLinha(linhaLetraA)
	if runa != 'A' {
		t.Errorf("Esperava 'A', veio %q", runa) // ➊
	}
	const nomeA = "LATIN CAPITAL LETTER A" // ➋
	if nome != nomeA {
		t.Errorf("Esperava %q, veio %q", nomeA, nome) // ➌
	}
}
```

➊ Para ser consistente com o próximo teste, mudei o código de formatação aqui de `%c` para `%q`, assim a runa vai aparecer entre aspas simples se o teste falhar.

➋ Criamos uma constante só para esse teste.

➌ O código de formatação `%q` também serve para exibir uma string entre aspas duplas.

Rodamos o teste para ver o que acontece:

```bash
$ go test
--- FAIL: TestAnalisarLinha (0.00s)
	runefinder_test.go:14: Esperava "LATIN CAPITAL LETTER A", veio "LETRA A"
FAIL
exit status 1
FAIL	github.com/ThoughtWorksInc/runas	0.011s
```

Note a mensagem indicando que na linha 14 um teste falhou. O texto da mensagem é o que passamos para o método `t.Errorf`.

Agora vamos codar a função `AnalisarLinha` de verdade no arquivo `runefinder.go`, assim:

```go
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
```

Explicando:

➊ Para importar dois ou mais pacotes, essa é a sintaxe utilizada.

➋ Todo identificaor importado de outro pacote é usado assim: `pacote.Identificador` (na verdade, é possível importar pacotes de outra forma, mas essa é a forma mais comum e mais recomendada). A função `strings.Split` recebe uma `string` para quebrar e outra `string` com o separador, e devolve uma fatia (_slice_) de strings, que é como um `array` de tamanho variável.

➌ A função `strconv.ParseInt` converte de `string` para `int64`. Ela recebe uma `string` (no caso, o item 0 da fatia `campos`), uma base (16) e o número de bits que se espera encontrar no inteiro resultante (32). O resultado é um `int64` e um objeto do tipo `error`, que nós vamos ignorar neste caso porque vamos assumir que as pessoas do Unicode sabem escrever números hexadecimais.

➍ Os valores devolvidos são o `código` convertido de `int64` para `rune`, e o segundo campo, que contém o nome do caractere.

Podemos rodar o teste para conferir o trabalho até aqui (com a opção `-v` para ver mais detalhes):

```bash
$ go test -v
=== RUN TestAnalisarLinha
--- PASS: TestAnalisarLinha (0.00s)
PASS
ok  	github.com/ThoughtWorksInc/runas	0.012s
```

Assim completamos o passo 2. Hora e ler o [Passo 3](passo-03).
