# Runas, passo 2: primeira função completa

Agora vamos terminar de implementar a função `AnalisarLinha`, guiados por testes.

Alteramos `runefinder_test.go` para verificar o nome devolvido por `AnalisarLinha`:

```go
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
```

<1> Para ser consistente com o próximo teste, mudei o código de formatação aqui de `%c` para `%q`, assim a runa vai aparecer entre aspas simples se o teste falhar.

<2> Criamos uma constante só para esse teste.

<3> O código de formatação `%q` também serve para exibir uma string entre aspas duplas.

Rodamos o teste para ver o que acontece:

```bash
$ go test
--- FAIL: TestAnalisarLinha (0.00s)
	runefinder_test.go:14: Esperava "LATIN CAPITAL LETTER A", veio "LETRA A"
FAIL
exit status 1
FAIL	github.com/labgo/runas-passo-a-passo	0.011s
```

Note a mensagem indicando que na linha 14 um teste falhou. O texto da mensagem é o que passamos para o métido `t.Errorf`.

Agora vamos codar a função `AnalisarLinha` de verdade, assim:

```go
package runefinder

import ( // <1>
	"strconv"
	"strings"
)

func AnalisarLinha(ucdLine string) (rune, string) {
	campos := strings.Split(ucdLine, ";")           // <2>
	código, _ := strconv.ParseInt(campos[0], 16, 32) // <3>
	return rune(código), campos[1]                  // <4>
}
```

Explicando:

<1> Para importar dois ou mais pacotes, essa é a sintaxe utilizada.

<2> Todo identificaor de outro pacote é usado assim: `pacote.Identificador` (na verdade, é possível importar identificadores de outra forma, mas essa é a forma mais comum e mais recomendada). A função `strings.Split` recebe uma `string` para quebrar e outra `string` com o separador, e devolve uma fatia (_slice_) de strings, que é como um `array` de tamanho variável.

<3> A função `strconv.ParseInt` converte de `string` para `int64`. Ela recebe uma `string` (no caso, o item 0 da fatia `campos`), uma base (16) e o número de bits que se espera encontrar no inteiro resultante (32). O resultado é um `int64` e um objeto do tipo `error`, que nós vamos ignorar neste caso porque vamos assumir que as pessoas do Unicode sabem escrever números hexadecimais.

Assim completamos o passo 2. Hora de mudar para o _branch_ `passo-03` e ler o arquivo `passo-03.md`.
